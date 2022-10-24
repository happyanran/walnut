package model

import (
	"strconv"
	"time"
)

type Dir struct {
	ID        int       `gorm:"type:integer NOT NULL PRIMARY KEY AUTOINCREMENT;" json:"ID"`
	PID       int       `gorm:"not null;uniqueIndex:idx_pid_name,priority:1;" json:"PID"`
	Name      string    `gorm:"not null;uniqueIndex:idx_pid_name,priority:2;" json:"name"`
	Path      string    `gorm:"not null;index:idx_path;" json:"path"`
	Note      string    `gorm:"" json:"note"`
	CreatedAt time.Time `gorm:"autoCreateTime;" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;" json:"updatedAt"`
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

func (m *Dir) DirDelete() error {
	return svcCtx.SqlDB.Delete(m).Error
}

func (m *Dir) DirFindByID() error {
	return svcCtx.SqlDB.First(m).Error
}

// 有问题
func (m *Dir) DirIterateDel() error {
	rows, err := svcCtx.SqlDB.Model(&Dir{}).Where("id = ? or path LIKE '?'", m.ID, m.Path+strconv.Itoa(m.ID)+","+"%").Rows()
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var dir Dir

		svcCtx.SqlDB.ScanRows(rows, &dir)

		if err := dir.DirDelete(); err != nil {
			return err
		}

		file := File{DirID: dir.ID}
		if err := file.FileDeleteByDirID(); err != nil {
			return err
		}

	}

	return nil
}

func (m *Dir) DirMoveChilds(oldPath string) error {
	//select path,'qqq'||substr(path,LENGTH('qqq')+1) from walnut_dir where path like 'abc%'
	newPath := m.Path + strconv.Itoa(m.ID) + ","
	return svcCtx.SqlDB.Exec("UPDATE walnut_dir SET path = ?||substr(path, LENGTH(?)+1) WHERE path like ?",
		newPath, oldPath, oldPath+"%").Error
}

func (m *Dir) FindChildDirs(dirs *[]Dir) error {
	return svcCtx.SqlDB.Model(&Dir{}).Where("p_id = ?", m.ID).Find(dirs).Error
}
