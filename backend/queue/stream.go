package queue

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"xinhang-backend/cache"
	"xinhang-backend/models"

	"github.com/redis/go-redis/v9"
)

const (
	streamName    = "xinhang:applications"
	consumerGroup = "app-workers"
	consumerName  = "worker-1"
)

var enabled bool
var cancelFn context.CancelFunc

func ConnectQueue() {
	if cache.RDB == nil {
		log.Println("WARNING: Redis not available, async queue disabled")
		enabled = false
		return
	}

	ctx := context.Background()
	err := cache.RDB.XGroupCreateMkStream(ctx, streamName, consumerGroup, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		log.Printf("WARNING: Redis Stream setup failed (%v), async queue disabled", err)
		enabled = false
		return
	}

	enabled = true
	log.Printf("Redis Stream queue ready, stream=%s, group=%s", streamName, consumerGroup)
}

func IsEnabled() bool {
	return enabled
}

func PublishApplication(app *models.Application) error {
	if !enabled {
		return nil
	}

	data, err := json.Marshal(app)
	if err != nil {
		return err
	}

	return cache.RDB.XAdd(context.Background(), &redis.XAddArgs{
		Stream: streamName,
		Values: map[string]interface{}{
			"email": app.Email,
			"data":  string(data),
		},
	}).Err()
}

func StartConsumer(handler func(app *models.Application) error) {
	if !enabled {
		return
	}

	var ctx context.Context
	ctx, cancelFn = context.WithCancel(context.Background())

	go func() {
		log.Println("Redis Stream consumer started")
		for {
			select {
			case <-ctx.Done():
				log.Println("Redis Stream consumer stopped")
				return
			default:
			}

			entries, err := cache.RDB.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    consumerGroup,
				Consumer: consumerName,
				Streams:  []string{streamName, ">"},
				Count:    10,
				Block:    5 * time.Second,
			}).Result()

			if err != nil {
				if err != redis.Nil {
					log.Printf("Stream consumer error: %v", err)
					time.Sleep(time.Second)
				}
				continue
			}

			for _, stream := range entries {
				for _, msg := range stream.Messages {
					data, ok := msg.Values["data"].(string)
					if !ok {
						cache.RDB.XAck(ctx, streamName, consumerGroup, msg.ID)
						continue
					}

					var app models.Application
					if err := json.Unmarshal([]byte(data), &app); err != nil {
						log.Printf("Stream message parse error: %v", err)
						cache.RDB.XAck(ctx, streamName, consumerGroup, msg.ID)
						continue
					}

					if err := handler(&app); err != nil {
						log.Printf("Stream handler error: %v", err)
						continue
					}

					cache.RDB.XAck(ctx, streamName, consumerGroup, msg.ID)
				}
			}
		}
	}()
}

func Close() {
	if cancelFn != nil {
		cancelFn()
	}
}
