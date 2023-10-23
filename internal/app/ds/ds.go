package ds

import "time"

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
	UserID       uint   `gorm:"primaryKey"`
	Name         string `gorm:"size:100"`
	PhoneNumber  string `gorm:"size:15"`
	EmailAddress string `gorm:"unique;size:100"`
	Password     string `gorm:"size:100"`
	Role         string `gorm:"size:100"`
}
