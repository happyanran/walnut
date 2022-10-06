package common

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewSqlite(s SqliteConf, l *logrus.Logger) (*gorm.DB, error) {
	newLogger := logger.New(
		l, // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	db, err := gorm.Open(sqlite.Open(s.path), &gorm.Config{
		Logger:                                   newLogger,
		CreateBatchSize:                          1000,
		DisableForeignKeyConstraintWhenMigrating: true,
		AllowGlobalUpdate:                        false,
		SkipDefaultTransaction:                   true,
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}
