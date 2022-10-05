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

		var d = &User{
			ID:       1,
			Username: "admin",
			Password: svcCtx.Utilw.PwdEnrypt("admin"),
		}

		if err := d.UserFindByID(); err == nil {
			return nil
		}

		d.UserCreate()
	}

	return nil
}
