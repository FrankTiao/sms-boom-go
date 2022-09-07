package utils

import (
	"log"
	"os"
	"strconv"
	"time"
)

func InitLog() {
	log.SetPrefix("[version:" + Version + "]")

	daily := time.Now().Format("2006-01-02")
	if _, err := DirExistsOrCreate(GetAppDataLogDir(daily)); err != nil {
		return
	}

	file := daily + "/" + strconv.FormatInt(time.Now().Unix(), 10) + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		return
	}

	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
}
