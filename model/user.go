package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int       `gorm:"type:integer NOT NULL PRIMARY KEY AUTOINCREMENT;"`
	Username  string    `gorm:"not null; unique;"`
	Password  string    `gorm:"not null;"`
	CreatedAt time.Time `gorm:"autoCreateTime;"`
}

func (User) TableName() string {
	return "walnut_user"
}

// NOTE Any zero value like 0, ”, false won’t be saved into the database for those fields defined default value,
// you might want to use pointer type or Scanner/Valuer to avoid this
func (u *User) UserCreate() error {
	return svcCtx.SqlDB.Create(u).Error
}

// Update with conditions and model value: db.Model(&user).Where("active = ?", true).Update("name", "hello")
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111 AND active=true;
// When update with struct, GORM will only update non-zero fields, you might want to use map to update attributes or use Select to specify fields to update
func (u *User) UserUpdate() error {
	if u.ID == 0 {
		return errors.New("主键必须非0值")
	}

	return svcCtx.SqlDB.Model(u).Omit("Username").Updates(u).Error
	//svcCtx.SqlDB.Save(u)
}

func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	if u.Username == "admin" {
		return errors.New("admin user not allowed to delete")
	}
	return
}

func (u *User) UserDelete() error {
	if u.ID == 0 {
		return errors.New("主键必须非0值")
	}

	return svcCtx.SqlDB.Delete(u).Error
}

func (u *User) UserFindByID() error {
	if u.ID == 0 {
		return errors.New("主键必须非0值")
	}

	return svcCtx.SqlDB.Omit("password").First(u).Error
}

// db.Where(&User{Name: "jinzhu"}, "name", "Age").Find(&users)
// SELECT * FROM users WHERE name = "jinzhu" AND age = 0;
func (u *User) UserFindByName() error {
	if u.Username == "" {
		return errors.New("Username必须非0值")
	}

	return svcCtx.SqlDB.Where(u, "username").Find(u).Error
}

// db.Order("age desc, name").Find(&users)
// SELECT * FROM users ORDER BY age desc, name;
// db.Limit(10).Offset(5).Find(&users)
// SELECT * FROM users OFFSET 5 LIMIT 10;
func (u *User) UserGetAll(users *[]User, limit, offset int) error {
	return svcCtx.SqlDB.Omit("password").Order("id").Limit(limit).Offset(offset).Find(users).Error
}

func (u *User) UserCount(cnt *int64) error {
	return svcCtx.SqlDB.Model(&User{}).Count(cnt).Error
}
