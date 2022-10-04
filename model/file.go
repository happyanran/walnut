package model

import "time"

type File struct {
	ID        int    `gorm:"type:integer NOT NULL PRIMARY KEY AUTOINCREMENT;"`
	DirID     int    `gorm:"not null;"`
	Name      string `gorm:"not null;"`
	ExtType   string `gorm:"not null;"`
	Hash      string
	Note      string
	Size      int       `gorm:"not null;"`
	CreatedAt time.Time `gorm:"autoCreateTime;"`
}

func (File) TableName() string {
	return "walnut_file"
}
