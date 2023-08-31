package logger

import (
	"log"
	"os"
)

// var ConsensusTimeLogger = initConsensusTimeLogger()
var BandWidthLogger = initDownLoadBandWidthLogger()
var UBandWidthLogger = initUpLoadBandWidthLoggger()

// func initConsensusTimeLogger() *log.Logger {
// 	logFile, _ := os.Create("ConsensusTime.log")
// 	logger := log.New(logFile, "[time]", log.Ltime|log.Lshortfile|log.LUTC) // 将文件设置为loger作为输出
// 	return logger
// }

func initDownLoadBandWidthLogger() *log.Logger {
	logFile, _ := os.Create("/Users/xiading/Library/Mobile Documents/com~apple~CloudDocs/学习/中山大学/论文代码/validator/DownLoadBandWidth.log")
	logger := log.New(logFile, "[bandwidth]", log.Ltime|log.Lshortfile|log.LUTC) // 将文件设置为loger作为输出
	return logger
}

func initUpLoadBandWidthLoggger() *log.Logger {
	logFile, _ := os.Create("/Users/xiading/Library/Mobile Documents/com~apple~CloudDocs/学习/中山大学/论文代码/validator/UpLoadBandWidth.log")
	logger := log.New(logFile, "[bandwidth]", log.Ltime|log.Lshortfile|log.LUTC) // 将文件设置为loger作为输出
	return logger
}
