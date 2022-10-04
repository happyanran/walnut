package model

import "time"

type Dir struct {
	ID        int    `gorm:"type:integer NOT NULL PRIMARY KEY AUTOINCREMENT;"`
	Pid       int    `gorm:"not null;"`
	Name      string `gorm:"not null;"`
	Note      string
	CreatedAt time.Time `gorm:"autoCreateTime;"`
}

func (Dir) TableName() string {
	return "walnut_dir"
}
