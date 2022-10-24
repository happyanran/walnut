package common

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Cfg    *Config
	Log    *logrus.Logger
	FileOp *FileOp
	SqlDB  *gorm.DB
	ZhVal  *Vlidate
	Jwtw   *Jwts
	Utilw  *Utils
}

func NewServiceContext(c *Config) *ServiceContext {
	log := NewLog(c.LogConf)

	fileOp, err := NewFileOp(c.ServerConf, log)
	if err != nil {
		log.Fatal("Dir init failed: ", err)
	}

	db, err := NewSqlite(c.SqliteConf, log)
	if err != nil {
		log.Fatal("Sqlite init failed: ", err)
	}

	val, err := NewValidate()
	if err != nil {
		log.Error("Validate init failed: ", err)
	}

	jwtw := NewJwtw(c.JwtConf)

	utilw := NewUtil()

	return &ServiceContext{
		Cfg:    c,
		Log:    log,
		FileOp: fileOp,
		SqlDB:  db,
		ZhVal:  val,
		Jwtw:   jwtw,
		Utilw:  utilw,
	}
}
