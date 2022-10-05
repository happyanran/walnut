package common

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	ServerConf
	LogConf
	SqliteConf
	JwtConf
}

type ServerConf struct {
	GinMode      string
	Addr         string
	MigrateTable bool
	Data         string
}

type LogConf struct {
	Level uint32
}

type SqliteConf struct {
	path string
}

type JwtConf struct {
	key        string
	expireHour int
}

func LoadConfig() *Config {
	v := viper.New()

	pflag.String("config", "./walnut.yaml", "config file path.")
	pflag.String("server.ginmode", "release", "The Gin mode.")
	pflag.String("server.addr", "0.0.0.0:8080", "The address to listen on for HTTP requests.")
	pflag.Bool("server.migratetable", true, "Auto migrate table.")
	pflag.String("server.data", "./data", "Root Dir.")
	pflag.String("log.level", "info", "log level: error, warn, info.")
	pflag.String("sqlite.path", "./walnut.db", "Sqlite db file path.")
	pflag.String("jwt.key", "abcdefg", "Jwt key.")
	pflag.Int("jwt.expirehour", 24, "Jwt key expire hours.")

	v.BindPFlags(pflag.CommandLine)
	pflag.Parse()

	v.SetConfigFile(v.GetString("config"))
	err := v.ReadInConfig()
	if err != nil {
		fmt.Println("Warning: Read config file failed: ", err.Error())
	}

	var level uint32
	switch v.GetString("log.level") {
	case "error":
		level = 2
	case "warn", "warning":
		level = 3
	case "info":
		level = 4
	}

	return &Config{
		ServerConf: ServerConf{
			GinMode:      v.GetString("server.ginmode"),
			Addr:         v.GetString("server.addr"),
			MigrateTable: v.GetBool("server.migratetable"),
			Data:         v.GetString("server.data"),
		},
		LogConf: LogConf{
			Level: level,
		},
		SqliteConf: SqliteConf{
			path: v.GetString("sqlite.path"),
		},
		JwtConf: JwtConf{
			key:        v.GetString("jwt.key"),
			expireHour: v.GetInt("jwt.expirehour"),
		},
	}
}
