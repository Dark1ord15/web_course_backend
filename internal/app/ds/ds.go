package ds

import (
	"Road_services/internal/app/role"
	"time"
)

type Road struct {
	Roadid          uint   `gorm:"primaryKey"`
	Name            string `gorm:"size:100"`
	Trustmanagment  int
	Length          int
	Paidlength      int
	Category        string `gorm:"size:50"`
	Numberofstripes string `gorm:"size:10"`
	Speed           int
	Price           int
	Image           string `gorm:"size:255"`
	Statusroad      string `gorm:"size:10"`
	Startofsection  int
	Endofsection    int
}

type Travelrequest struct {
	Travelrequestid uint `gorm:"primaryKey"`
	Userid          uint
	Requeststatus   string    `gorm:"size:100"`
	Creationdate    time.Time `gorm:"default:NULL"`
	Formationdate   time.Time `gorm:"default:NULL"`
	Completiondate  time.Time `gorm:"default:NULL"`
	Moderatorid     uint      `gorm:"default:NULL"`
}

type Travelrequestroad struct {
	Travelrequestroadid uint `gorm:"primaryKey"`
	Travelrequestid     uint
	Roadid              uint
}

type User struct {
	Id          uint      `gorm:"primarykey"`
	Name        string    `json:"name"`
	Login       string    `json:"login"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	Role        role.Role `sql:"type:string;"`
	Password    string    `gorm:"size:60"`
}
