package setting

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type LogHandler struct {
}

var Log_Ins *LogHandler
var Log_once sync.Once

func GetLogManager() *LogHandler {
	Log_once.Do(func() {
		Log_Ins = &LogHandler{}
	})
	return Log_Ins
}

func (l *LogHandler) SetLogFile(Port string) {
	// 현재시간
	startDate := time.Now().Format("2006-01-02")
	// log 폴더 위치
	logFolderPath := "/dipnas/DIPServer/ServerLog/" + GetServerType(Port)

	// log 파일 경로
	logFilePath := fmt.Sprintf("%slogFile-%s.log", logFolderPath, startDate)

	// log 폴더가 없을 경우 log 폴더 생성
	if _, err := os.Stat(logFolderPath); os.IsNotExist(err) {
		os.MkdirAll(logFolderPath, 0777)
	}

	// log 파일이 없을 경우 log 파일 생성
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		os.Create(logFilePath)
	}
	// log 파일 열기
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	// log 패키지를 활요하여 작성할 경우 log 파일에 작성되도록 설정
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.Println("LogStored In", logFolderPath)
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "err"
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "err"
}

func GetServerType(port string) string {
	switch GetLocalIP() + port { // 서버 실행 IP 마다 다른 폴더에서 로그저장
	//PPRK 내부 서버
	case "192.168.0.9:8001":
		return "PPRKContent/"
	case "192.168.0.9:8099":
		return "PPRKTEST/"
	case "192.168.0.9:8002":
		return "PPRKScreen/"
	case "192.168.0.9:8005":
		return "PPRKVoice/"
	case "192.168.0.9:3000":
		return "PPRKLoginLauncher/"

		//원효로 서버
	case "10.5.147.193:8000":
		return "Dev/"
	case "10.5.147.184:8000":
		return "Content1/"
	case "10.5.147.192:8000":
		return "Content2/"
	case "10.5.147.200:3000":
		return "LoginLauncher1/"
	case "10.5.147.107:3000":
		return "LoginLauncher2/"
	case "10.5.147.200:8000":
		return "Screen1/"
	case "10.5.147.107:8000":
		return "Screen2/"
	case "10.5.147.169:8000":
		return "Voice1/"
	case "10.5.147.131:8000":
		return "Voice2/"
	case "10.5.147.184:4401":
		return "File1/"
	case "10.5.147.192:4401":
		return "File2/"
	}
	return "err/"
}
