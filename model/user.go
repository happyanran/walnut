package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int       `gorm:"type:integer NOT NULL PRIMARY KEY AUTOINCREMENT;" json:"ID"`
	UserName  string    `gorm:"not null; uniqueIndex:idx_username;" json:"userName"`
	Password  string    `gorm:"not null;" json:"-"`
	NickName  string    `gorm:"" json:"nickName"`
	CreatedAt time.Time `gorm:"autoCreateTime;" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;" json:"updatedAt"`
	// Role Todo
}

func (User) TableName() string {
	return "walnut_user"
}

// NOTE Any zero value like 0, ”, false won’t be saved into the database for those fields defined default value,
// you might want to use pointer type or Scanner/Valuer to avoid this
func (m *User) UserCreate() error {
	return svcCtx.SqlDB.Create(m).Error
}

// Update with conditions and model value: db.Model(&user).Where("active = ?", true).Update("name", "hello")
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111 AND active=true;
// When update with struct, GORM will only update non-zero fields, you might want to use map to update attributes or use Select to specify fields to update
func (m *User) UserUpdate() error {
	return svcCtx.SqlDB.Model(m).Omit("user_name").Updates(m).Error
	//svcCtx.SqlDB.Save(u)
}

func (m *User) BeforeDelete(tx *gorm.DB) (err error) {
	if m.ID == 1 {
		return errors.New("admin user not allowed to delete")
	}
	return
}

func (m *User) UserDelete() error {
	return svcCtx.SqlDB.Delete(m).Error
}

func (m *User) UserFindByID() error {
	return svcCtx.SqlDB.Omit("password").First(m).Error
}

// db.Where(&User{Name: "jinzhu"}, "name", "Age").Find(&users)
// SELECT * FROM users WHERE name = "jinzhu" AND age = 0;
func (m *User) UserFindByName() error {
	return svcCtx.SqlDB.Where(m, "user_name").Find(m).Error
}

// db.Order("age desc, name").Find(&users)
// SELECT * FROM users ORDER BY age desc, name;
// db.Limit(10).Offset(5).Find(&users)
// SELECT * FROM users OFFSET 5 LIMIT 10;
func (m *User) UserGetAll(users *[]User, limit, offset int) error {
	return svcCtx.SqlDB.Omit("password").Order("id").Limit(limit).Offset(offset).Find(users).Error
}

func (m *User) UserCount(cnt *int64) error {
	return svcCtx.SqlDB.Model(&User{}).Count(cnt).Error
}
