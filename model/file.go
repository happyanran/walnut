package model

import (
	"time"
)

type File struct {
	ID        int       `gorm:"type:integer NOT NULL PRIMARY KEY AUTOINCREMENT;"`
	DirID     int       `gorm:"not null;uniqueIndex:idx_dirid_name,priority:1;"`
	Name      string    `gorm:"not null;uniqueIndex:idx_dirid_name,priority:2;"`
	ExtType   string    `gorm:"not null;"`
	Note      string    `gorm:""`
	Hash      string    `gorm:""`
	Size      int       `gorm:"not null;"`
	CreatedAt time.Time `gorm:"autoCreateTime;"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;"`
}

func (File) TableName() string {
	return "walnut_file"
}

func (m *File) FileCreate() error {
	return svcCtx.SqlDB.Create(m).Error
}

func (m *File) FileUpdate() error {
	return svcCtx.SqlDB.Save(m).Error
}

func (m *File) FileDelete() error {
	return svcCtx.SqlDB.Delete(m).Error
}

func (m *File) FileFindByID() error {
	return svcCtx.SqlDB.First(m).Error
}

func (m *File) FileFindByDirID(files *[]File) error {
	return svcCtx.SqlDB.Where(m, "Dir_ID").Find(files).Error
}
