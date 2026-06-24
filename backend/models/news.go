package models

import "time"

type News struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"size:300;not null"`
	TitleEn     string    `json:"titleEn" gorm:"size:300"`
	Summary     string    `json:"summary" gorm:"size:1000"`
	SummaryEn   string    `json:"summaryEn" gorm:"size:1000"`
	Content     string    `json:"content" gorm:"type:text"`
	ContentEn   string    `json:"contentEn" gorm:"type:text"`
	Category    string    `json:"category" gorm:"size:50;index;not null;default:news"`
	CoverImage  string    `json:"coverImage" gorm:"size:500"`
	AuthorID    *uint     `json:"authorId" gorm:"index"`
	Published   bool      `json:"published" gorm:"default:false;index"`
	PublishedAt *time.Time `json:"publishedAt"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    *uint     `json:"userId" gorm:"index"`
	Title     string    `json:"title" gorm:"size:300;not null"`
	Content   string    `json:"content" gorm:"type:text"`
	Type      string    `json:"type" gorm:"size:50;index;not null;default:system"`
	IsRead    bool      `json:"isRead" gorm:"default:false;index"`
	RefType   string    `json:"refType" gorm:"size:50"`
	RefID     *uint     `json:"refId"`
	CreatedAt time.Time `json:"createdAt"`
}
