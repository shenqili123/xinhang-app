package models

import "time"

type Application struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        *uint     `json:"userId" gorm:"index"`
	PermitNo      string    `json:"permitNo" gorm:"size:30;uniqueIndex"`
	StudentName   string    `json:"studentName" gorm:"size:100;not null"`
	BirthDate     string    `json:"birthDate" gorm:"size:20"`
	Gender        string    `json:"gender" gorm:"size:20;not null"`
	Division      string    `json:"division" gorm:"size:30"`
	Grade         int       `json:"grade" gorm:"not null"`
	IDNumber      string    `json:"idNumber" gorm:"size:30"`
	Photo         string    `json:"photo" gorm:"size:500"`
	ExamRoom      string    `json:"examRoom" gorm:"size:10"`
	SeatNumber    string    `json:"seatNumber" gorm:"size:10"`
	BoardingNeed  string    `json:"boardingNeed" gorm:"size:30"`
	ParentName    string    `json:"parentName" gorm:"size:100;not null"`
	Phone         string    `json:"phone" gorm:"size:20;not null;index"`
	Relationship  string    `json:"relationship" gorm:"size:30"`
	Email         string    `json:"email" gorm:"size:200"`
	CurrentSchool string    `json:"currentSchool" gorm:"size:200"`
	Track         string    `json:"track" gorm:"size:100"`
	VisitDate     string    `json:"visitDate" gorm:"size:20"`
	Notes         string    `json:"notes" gorm:"type:text"`
	Status        string    `json:"status" gorm:"size:20;default:pending"`
	CreatedAt     time.Time `json:"createdAt"`
}

type ApplyRequest struct {
	StudentName   string `json:"studentName" binding:"required,max=100"`
	BirthDate     string `json:"birthDate" binding:"max=20"`
	Gender        string `json:"gender" binding:"required,max=20"`
	Division      string `json:"division" binding:"max=30"`
	Grade         int    `json:"grade" binding:"min=0,max=12"`
	IDNumber      string `json:"idNumber" binding:"required,max=30"`
	Photo         string `json:"photo" binding:"required,max=500"`
	BoardingNeed  string `json:"boardingNeed" binding:"required,max=30"`
	ParentName    string `json:"parentName" binding:"required,max=100"`
	Phone         string `json:"phone" binding:"required,min=8,max=20"`
	Relationship  string `json:"relationship" binding:"required,max=30"`
	Email         string `json:"email" binding:"required,max=200"`
	CurrentSchool string `json:"currentSchool" binding:"required,max=200"`
	Track         string `json:"track" binding:"required,max=100"`
	VisitDate     string `json:"visitDate" binding:"required,max=20"`
	Notes         string `json:"notes" binding:"max=2000"`
}
