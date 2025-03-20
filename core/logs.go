package core

import (
	"fmt"
	"log"
	"time"
)

const (
	logPath = "logs/guide.log"
)

var (
	mlog = Wlogs{}
)

type Logs interface {
	Info(string)
	Error(string)
	Warn(string)
}

type Wlogs struct {
}

func (w Wlogs) Info(l string) {
	now := time.Now().Format("2006/01/02 - 15:04:05")
	mlog := fmt.Sprintf("[GUIDE] %v INFO %s", now, l)
	if err := WriteFile(logPath, mlog); err != nil {
		log.Println(err.Error())
	}
}

func (w Wlogs) Error(l string) {
	now := time.Now().Format("2006/01/02 - 15:04:05")
	mlog := fmt.Sprintf("[GUIDE] %v ERROR %s", now, l)
	if err := WriteFile(logPath, mlog); err != nil {
		log.Println(err.Error())
	}
}

func (w Wlogs) Warn(l string) {
	now := time.Now().Format("2006/01/02 - 15:04:05")
	mlog := fmt.Sprintf("[GUIDE] %v WARN %s", now, l)
	if err := WriteFile(logPath, mlog); err != nil {
		log.Println(err.Error())
	}
}
