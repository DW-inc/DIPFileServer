package setting

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type SettingHandler struct {
	ServerType int // 0: 나, 1: 원효로1번서버, 2: 원효로2번서버
	Port       string
	NasPath    string
	CTSAddress string
	LogPath    string
	DB         string
}

var St_Ins *SettingHandler
var St_once sync.Once

func GetStManager() *SettingHandler {
	St_once.Do(func() {
		St_Ins = &SettingHandler{}
	})
	return St_Ins
}

func (st *SettingHandler) Init() {
	st.ServerType = 0 // 0: 나, 1: 원효로1번서버, 2: 원효로2번서버

	err := godotenv.Load("./setting/process.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	switch st.ServerType {
	case 0:
		st.Port = ":8009"
		st.NasPath = "../Server/Storage/nas/"
		st.LogPath = "/data/DIPServerLog/FileServer/"
		st.CTSAddress = "192.168.0.9:8001"
		st.DB = os.Getenv("CONTENT_DB0")
	case 1:
		st.Port = ":4401"
		st.NasPath = "/dipnas/DIPServer/Storage/"
		st.LogPath = "/data/DIPServerLog/FileServer1/"
		st.CTSAddress = "10.5.147.88:8000"
		st.DB = os.Getenv("CONTENT_DB1")
	case 2:
		st.Port = ":4401"
		st.NasPath = "/dipnas/DIPServer/Storage/"
		st.LogPath = "/data/DIPServerLog/FileServer2/"
		st.CTSAddress = "10.5.147.88:8000"
		st.DB = os.Getenv("CONTENT_DB1")
	}
}
