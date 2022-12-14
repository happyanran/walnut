package model

import (
	"github.com/happyanran/walnut/common"
)

var svcCtx *common.ServiceContext

func Init(s *common.ServiceContext) error {
	svcCtx = s

	if svcCtx.Cfg.ServerConf.MigrateTable {
		err := svcCtx.SqlDB.AutoMigrate(
			User{},
			Dir{},
			File{},
		)
		if err != nil {
			return err
		}

		u := &User{
			ID:       1,
			UserName: "admin",
			Password: svcCtx.Utilw.PwdEnrypt("admin"),
			NickName: "admin",
		}

		if err := u.UserFindByID(); err == nil {
			return nil
		}

		u.UserCreate()

		d := &Dir{
			ID:   1,
			PID:  0,
			Name: "小核桃",
		}

		if err := d.DirFindByID(); err == nil {
			return nil
		}

		d.DirCreate()
	}

	return nil
}
