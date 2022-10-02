package common

import (
	"github.com/sirupsen/logrus"
)

func NewLog(l LogConf) *logrus.Logger {
	log := logrus.New()

	log.SetLevel(logrus.Level(l.Level))
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
	//log.WithFields(log.Fields{  "event": event,  "topic": topic,  "key": key,}).Fatal("Failed to send event")

	return log
}
