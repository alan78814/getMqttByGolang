package service

import "github.com/sirupsen/logrus"

// 初始化 logger
var Logger = logrus.New()

func Init() {
	Logger.SetLevel(logrus.InfoLevel)
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

}
