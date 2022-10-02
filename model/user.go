package model

import "time"

type User struct {
	ID        uint32    `gorm:"type:integer NOT NULL PRIMARY KEY AUTOINCREMENT;"`
	Username  string    `gorm:"not null; unique;"`
	Password  string    `gorm:"not null;"`
	CreatedAt time.Time `gorm:"autoCreateTime;"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;"`
}

func (User) TableName() string {
	return "walnut_file"
}

func (u *User) UserCreate() error {
	return svcCtx.SqlDB.Create(u).Error
}

func (u *User) UserUpdate() error {
	return svcCtx.SqlDB.Save(u).Error
}

func (u *User) UserDelete() error {
	return svcCtx.SqlDB.Delete(u).Error
}

func (u *User) UserFind() error {
	return svcCtx.SqlDB.First(u).Error
}

func (u *User) UserFindByName() error {
	return svcCtx.SqlDB.Where("username = ?", u.Username).First(u).Error
}

func (u *User) UserFindSignin() error {
	return svcCtx.SqlDB.Where("username = ? AND password = ?", u.Username, u.Password).First(u).Error
}
