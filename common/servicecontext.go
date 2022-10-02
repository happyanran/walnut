package common

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Cfg   *Config
	Log   *logrus.Logger
	SqlDB *gorm.DB
	ZhVal *Vlidate
}

func NewServiceContext(c *Config) *ServiceContext {
	log := NewLog(c.LogConf)
	
	db,err := NewSqlite(c.SqliteConf, log)
	if err != nil {
		log.Error("Sqlite init failed: ", err)
	}

	val, err := NewValidate()
	if err != nil {
		log.Error("Validate init failed: ", err)
	}

	return &ServiceContext{
		Cfg:   c,
		Log:   log,
		SqlDB: db,
		ZhVal: val,
	}
}
