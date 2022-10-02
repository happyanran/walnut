package model

import "time"

type File struct {
	ID        uint32 `gorm:"type:integer NOT NULL PRIMARY KEY AUTOINCREMENT;"`
	DirID     int    `gorm:"not null;"`
	Name      string `gorm:"not null;"`
	ExtType   string `gorm:"not null;"`
	Hash      string
	Note      string
	size      int       `gorm:"not null;"`
	CreatedAt time.Time `gorm:"autoCreateTime;"`
}

func (File) TableName() string {
	return "walnut_file"
}
