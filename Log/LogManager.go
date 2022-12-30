package setting

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/DW-inc/FileServer/setting"
)

type LogHandler struct {
	Day           string
	lock          sync.Mutex
	logFile       *os.File
	logNum        int
	logFolderPath string
	logFilePath   string
	MAXLOGSIZE    int64
}

var Log_Ins *LogHandler
var Log_once sync.Once

func GetLogManager() *LogHandler {
	Log_once.Do(func() {
		Log_Ins = &LogHandler{}
	})
	return Log_Ins
}

func (l *LogHandler) SetLogFile() {
	l.MAXLOGSIZE = 5 * 1024 * 1024
	// 현재시간
	startDate := time.Now().Format("2006-01-02")
	l.Day = startDate
	l.logNum = 1
	// log 폴더 위치
	l.logFolderPath = setting.GetStManager().LogPath
	// log 폴더가 없을 경우 log 폴더 생성
	if _, err := os.Stat(l.logFolderPath); os.IsNotExist(err) {
		os.MkdirAll(l.logFolderPath, 0777)
	}

	// log 파일 결정하기
	l.logFilePath = fmt.Sprintf("%slogFile-%s_%d.log", l.logFolderPath, startDate, l.logNum)
	for {
		if tempfile, err := os.Stat(l.logFilePath); os.IsNotExist(err) {
			os.Create(l.logFilePath)
			break
		} else {
			log.Println(tempfile.Size())
			if tempfile.Size() > l.MAXLOGSIZE {
				l.logNum += 1
				l.logFilePath = fmt.Sprintf("%slogFile-%s_%d.log", l.logFolderPath, startDate, l.logNum)
			} else {
				break
			}
		}
	}
	l.SetLog()
	go l.LogRoutine()
}
func (l *LogHandler) SetLog() {
	// log 파일 열기
	logFile, err := os.OpenFile(l.logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	// log 패키지를 활요하여 작성할 경우 log 파일에 작성되도록 설정
	mw := io.MultiWriter(os.Stdout, logFile)
	l.logFile = logFile
	log.SetOutput(mw)
	log.Println("LogStored In", l.logFilePath)
}

func (l *LogHandler) LogRoutine() {
	for {
		time.Sleep(time.Second * 1)

		// 현재시간
		startDate := time.Now().Format("2006-01-02")

		if startDate != l.Day {
			l.Day = startDate
			l.logNum = 1
			l.logFilePath = fmt.Sprintf("%slogFile-%s_1.log", l.logFolderPath, startDate)
			if err := l.LogRotate(); err != nil {
				log.Println(err)
			}
		} else {
			tempfile, _ := os.Stat(l.logFilePath)
			if tempfile.Size() > l.MAXLOGSIZE {
				l.logNum += 1
				l.logFilePath = fmt.Sprintf("%slogFile-%s_%d.log", l.logFolderPath, startDate, l.logNum)
				if err := l.LogRotate(); err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func (l *LogHandler) LogRotate() error {
	var err error

	l.lock.Lock()
	defer l.lock.Unlock()

	if l.logFile != nil {
		err = l.logFile.Close()
		l.logFile = nil
		if err != nil {
			return err
		}
	}

	if err == nil {
		_, err = os.Create(l.logFilePath)
		if err != nil {
			return err
		}
		l.SetLog()
	}
	return nil
}
