package ds

import "time"

type Road struct {
	RoadID          uint   `gorm:"primaryKey"`
	Name            string `gorm:"size:100"`
	TrustManagement int
	Length          int
	PaidLength      int
	Category        string `gorm:"size:50"`
	NumberOfStripes string `gorm:"size:10"`
	Speed           int
	Price           int
	Image           string `gorm:"size:255"`
	StatusRoad      string `gorm:"size:10"`
}

type TravelRequest struct {
	TravelRequestID uint `gorm:"primaryKey"`
	UserID          uint
	RequestStatus   string `gorm:"size:100"`
	CreationDate    time.Time
	FormationDate   time.Time
	CompletionDate  time.Time
	ModeratorID     uint
}

type TravelRequestRoad struct {
	TravelRequestRoadID uint `gorm:"primaryKey"`
	TravelRequestID     uint
	RoadID              uint
}

type User struct {
	UserID       uint   `gorm:"primaryKey"`
	Name         string `gorm:"size:100"`
	PhoneNumber  string `gorm:"size:15"`
	EmailAddress string `gorm:"unique;size:100"`
	Password     string `gorm:"size:100"`
	Role         string `gorm:"size:100"`
}
