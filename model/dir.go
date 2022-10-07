package model

import (
	"path/filepath"
	"strconv"
	"time"
)

type Dir struct {
	ID        int       `gorm:"type:integer NOT NULL PRIMARY KEY AUTOINCREMENT;"`
	PID       int       `gorm:"not null;uniqueIndex:idx_pid_name,priority:1;"`
	Path      string    `gorm:"not null;index:idx_path;"`
	Name      string    `gorm:"not null;uniqueIndex:idx_pid_name,priority:2;"`
	Note      string    `gorm:""`
	CreatedAt time.Time `gorm:"autoCreateTime;"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;"`
	Files     []File
}

func (Dir) TableName() string {
	return "walnut_dir"
}

func (m *Dir) DirCreate() error {
	return svcCtx.SqlDB.Create(m).Error
}

func (m *Dir) DirUpdate() error {
	return svcCtx.SqlDB.Save(m).Error
}

func (m *Dir) DirUpdateChild(newpath string) error {
	//select path,'qqq'||substr(path,LENGTH('qqq')+1) from walnut_dir where path like 'abc%'
	oldpath := filepath.Join(m.Path, strconv.Itoa(m.ID))
	return svcCtx.SqlDB.Exec("UPDATE walnut_dir SET path = ?||substr(path, LENGTH(?)+1) WHERE p_id = ? or path like ?",
		newpath, oldpath, m.ID, filepath.Join(oldpath, "%")).Error
}

func (m *Dir) DirDelete() error {
	return svcCtx.SqlDB.Delete(m).Error
}

func (m *Dir) DirDeleteChild() error {
	return svcCtx.SqlDB.Where("p_id = ? or path LIKE ?", m.ID, filepath.Join(m.Path, strconv.Itoa(m.ID), "%")).Delete(&Dir{}).Error
}

func (m *Dir) DirDelNested() error {
	if err := m.DirDeleteChild(); err != nil {
		return err
	}

	return m.DirDelete()
}

func (m *Dir) DirFindByID() error {
	return svcCtx.SqlDB.First(m).Error
}

func (m *Dir) DirFilesFindByID() error {
	return svcCtx.SqlDB.Preload("Files").First(m).Error
}

func (m *Dir) DirFindByPID(dirs *[]Dir) error {
	return svcCtx.SqlDB.Where(m, "p_id").Find(dirs).Error
}

func (m *Dir) DirNameCheckByPID(pid int, name string) (int64, error) {
	var count int64

	if err := svcCtx.SqlDB.Model(&Dir{}).Where("p_id = ? and name = ?", pid, name).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
