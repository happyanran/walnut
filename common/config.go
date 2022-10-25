package common

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	ServerConf
	LogConf
	SqliteConf
	JwtConf
	WebDavConf
	HttpsConf
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
	Path string
}

type JwtConf struct {
	Key        string
	ExpireHour int
}

type WebDavConf struct {
	Enable bool
	Addr   string
	Data   string
}

type HttpsConf struct {
	Enable   bool
	Certfile string
	Keyfile  string
}

func LoadConfig() *Config {
	v := viper.New()

	pflag.String("config", "./conf/walnut.yaml", "config file path.")
	pflag.String("server.ginmode", "release", "The Gin mode.")
	pflag.String("server.addr", "0.0.0.0:8081", "The address to listen on for HTTP requests.")
	pflag.Bool("server.migratetable", true, "Auto migrate table.")
	pflag.String("server.data", "./data/walnut", "Data Dir.")
	pflag.String("log.level", "info", "log level: error, warn, info.")
	pflag.String("sqlite.path", "./data/walnut.db", "Sqlite db file path.")
	pflag.String("jwt.key", "aabbccddeeffgg", "Jwt key.")
	pflag.Int("jwt.expirehour", 24, "Jwt key expire hours.")
	pflag.Bool("webdav.enable", true, "Enable WebDav.")
	pflag.String("webdav.addr", "0.0.0.0:8082", "The address to listen on for WebDav.")
	pflag.String("webdav.data", "./data/webdav", "Root dir.")
	pflag.Bool("https.enable", true, "Enable https.")
	pflag.String("https.certfile", "./conf/server.crt", "Cert file path.")
	pflag.String("https.keyfile", "./conf/server.key", "Key file path.")

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of Little Walnut\n")
		pflag.PrintDefaults()
	}

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
			Path: v.GetString("sqlite.path"),
		},
		JwtConf: JwtConf{
			Key:        v.GetString("jwt.key"),
			ExpireHour: v.GetInt("jwt.expirehour"),
		},
		WebDavConf: WebDavConf{
			Enable: v.GetBool("webdav.enable"),
			Addr:   v.GetString("webdav.addr"),
			Data:   v.GetString("webdav.data"),
		},
		HttpsConf: HttpsConf{
			Enable:   v.GetBool("https.enable"),
			Certfile: v.GetString("https.certfile"),
			Keyfile:  v.GetString("https.keyfile"),
		},
	}
}
