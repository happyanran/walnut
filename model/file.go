package model

import (
	"time"
)

type File struct {
	ID           int       `gorm:"type:integer NOT NULL PRIMARY KEY AUTOINCREMENT;" json:"ID"`
	DirID        int       `gorm:"not null;uniqueIndex:idx_dirid_name,priority:1;" json:"dirID"`
	Name         string    `gorm:"not null;uniqueIndex:idx_dirid_name,priority:2;" json:"name"`
	ExtType      string    `gorm:"not null;index:idx_exttype;" json:"extType"`
	OriginalName string    `gorm:"not null;" json:"originalName"`
	SmallImgName string    `gorm:"" json:"smallImgName"`
	LargeImgName string    `gorm:"" json:"largeImgName"`
	Note         string    `gorm:"" json:"note"`
	OriginalSize int64     `gorm:"not null;" json:"originalSize"`
	CreatedAt    time.Time `gorm:"autoCreateTime;" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime;" json:"updatedAt"`
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

func (m *File) FileDeleteByDirID() error {
	return svcCtx.SqlDB.Where("dir_id = ?", m.DirID).Delete(&File{}).Error
}

func (m *File) FileFindByDirID(files *[]File) error {
	return svcCtx.SqlDB.Where(m, "dir_id").Find(files).Error
}

func (m *File) FileFindByName() error {
	return svcCtx.SqlDB.Where("dir_id = ? and name = ?", m.DirID, m.Name).Find(m).Error
}
