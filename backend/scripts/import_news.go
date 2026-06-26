package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Article struct {
	ArticleID   string `json:"article_id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Source      string `json:"source"`
	PubDate     string `json:"pub_date"`
	ContentHTML string `json:"content_html"`
	CoverImage  string `json:"cover_image"`
	ColumnName  string `json:"column_name"`
	Keywords    string `json:"keywords"`
}

type News struct {
	ID          uint       `gorm:"primaryKey"`
	Title       string     `gorm:"size:300;not null"`
	TitleEn     string     `gorm:"size:300"`
	Summary     string     `gorm:"size:1000"`
	SummaryEn   string     `gorm:"size:1000"`
	Content     string     `gorm:"type:text"`
	ContentEn   string     `gorm:"type:text"`
	Category    string     `gorm:"size:50;index;not null;default:news"`
	CoverImage  string     `gorm:"size:500"`
	AuthorID    *uint      `gorm:"index"`
	AuthorName  string     `gorm:"size:100"`
	Source      string     `gorm:"size:200"`
	Keywords    string     `gorm:"size:500"`
	OldID       string     `gorm:"size:50;index"`
	Published   bool       `gorm:"default:false;index"`
	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var categoryMap = map[string]string{
	"校园动态":   "campus",
	"新闻动态":   "news",
	"名师风采":   "teachers",
	"学生活动":   "activities",
	"学子风采":   "highlights",
	"媒体聚焦":   "media",
	"教育科研":   "research",
	"教学教研":   "research",
	"教学活动":   "research",
	"科研动态":   "research",
	"德育天地":   "moral_education",
	"德育活动":   "moral_education",
	"All":    "news_en",
	"友好学校":   "partner_schools",
	"家校交流":   "home_school",
	"家校平台":   "home_school",
	"家委会建设":  "home_school",
	"父母大学":   "home_school",
	"党建工作":   "party",
	"招生专栏":   "admission",
	"国际部":    "international",
	"学生社团":   "activities",
	"心灵花园":   "moral_education",
	"通知公告":   "notice",
	"校园新闻":   "campus",
	"校务公开":   "public_affairs",
	"办学成果":   "achievements",
	"国际交流":   "international",
	"课程介绍":   "courses",
}

const imageBaseURL = "/uploads/migration/images/"

var styleRegex = regexp.MustCompile(`\s*style="[^"]*"`)
var styleRegex2 = regexp.MustCompile(`\s*style='[^']*'`)
var tagRegex = regexp.MustCompile(`<[^>]+>`)
var spaceRegex = regexp.MustCompile(`\s+`)

func cleanHTML(html string) string {
	html = styleRegex.ReplaceAllString(html, "")
	html = styleRegex2.ReplaceAllString(html, "")
	html = strings.ReplaceAll(html, `src="images/`, `src="`+imageBaseURL)
	html = strings.ReplaceAll(html, `src="media/`, `src="/uploads/migration/media/`)
	return html
}

func extractSummary(html string, maxLen int) string {
	text := tagRegex.ReplaceAllString(html, "")
	text = spaceRegex.ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	runes := []rune(text)
	if len(runes) > maxLen {
		text = string(runes[:maxLen]) + "..."
	}
	return text
}

func parseDate(dateStr string) time.Time {
	formats := []string{
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, dateStr); err == nil {
			return t
		}
	}
	return time.Now()
}

func main() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		envOr("DB_HOST", "xinhang-db-postgresql.ns-0h7fttt7.svc"),
		envOr("DB_USER", "postgres"),
		envOr("DB_PASSWORD", "7dl72vft"),
		envOr("DB_NAME", "xinhang"),
		envOr("DB_PORT", "5432"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	fmt.Println("Database connected")

	articlesPath := "/home/devbox/project/migration_output/articles.json"
	data, err := os.ReadFile(articlesPath)
	if err != nil {
		log.Fatalf("Failed to read articles.json: %v", err)
	}

	var articles []Article
	if err := json.Unmarshal(data, &articles); err != nil {
		log.Fatalf("Failed to parse articles.json: %v", err)
	}
	fmt.Printf("Loaded %d articles\n", len(articles))

	// Deduplicate (normalize punctuation for comparison)
	seen := make(map[string]bool)
	var unique []Article
	for _, a := range articles {
		title := strings.TrimSpace(a.Title)
		if title == "" {
			continue
		}
		normTitle := normalizeTitle(title)
		if seen[normTitle] {
			continue
		}
		seen[normTitle] = true
		unique = append(unique, a)
	}
	fmt.Printf("After dedup: %d unique (%d removed)\n", len(unique), len(articles)-len(unique))

	// Clear all news data before reimport
	db.Exec("DELETE FROM news")
	fmt.Println("Cleared previous migration data")

	imported := 0
	errors := 0
	for _, art := range unique {
		title := strings.TrimSpace(art.Title)
		content := cleanHTML(art.ContentHTML)
		summary := extractSummary(art.ContentHTML, 200)
		category := categoryMap[art.ColumnName]
		if category == "" {
			category = "other"
		}

		cover := art.CoverImage
		if cover != "" {
			cover = strings.ReplaceAll(cover, "images/", "")
			cover = imageBaseURL + cover
		}

		pubTime := parseDate(art.PubDate)

		news := News{
			Title:       title,
			Summary:     summary,
			Content:     content,
			Category:    category,
			CoverImage:  cover,
			AuthorName:  strings.TrimSpace(art.Author),
			Source:      strings.TrimSpace(art.Source),
			Keywords:    strings.TrimSpace(art.Keywords),
			OldID:       art.ArticleID,
			Published:   true,
			PublishedAt: &pubTime,
			CreatedAt:   pubTime,
			UpdatedAt:   pubTime,
		}

		if err := db.Create(&news).Error; err != nil {
			errors++
			fmt.Printf("  ERROR: %s - %v\n", title[:min(30, len(title))], err)
		} else {
			imported++
		}
	}

	fmt.Printf("\nDone!\n  Imported: %d\n  Errors: %d\n", imported, errors)
}

func normalizeTitle(s string) string {
	replacer := strings.NewReplacer(
		"\uff1a", ":", "\uff0c", ",", "\u3001", ",",
		"\uff1b", ";", "\uff01", "!", "\uff1f", "?",
		"\u3000", "", "\u201c", "\"", "\u201d", "\"",
		"\u2018", "'", "\u2019", "'",
		" ", "", "\t", "", "\n", "", "\r", "",
	)
	return replacer.Replace(s)
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
