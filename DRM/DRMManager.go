package drm

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

type DrmManager struct {
}

var DRM_Ins *DrmManager
var once sync.Once

func GetDrmManager() *DrmManager {
	once.Do(func() {
		DRM_Ins = &DrmManager{}
	})
	return DRM_Ins
}

func (fm *DrmManager) Init() {

}

func (fm *DrmManager) DRM_Encrypt(FileName string) bool {
	cmd := exec.Command("java", "-cp", "scsl.jar", "TestEnc.java", FileName)
	_, err := cmd.Output()
	if err != nil {
		log.Println("DRM_Err", err)
		return false
	}

	return fm.CheckFile("{ENC}"+FileName, "DRM Encrypt")
}

func (fm *DrmManager) DRM_Decrypt(FileName string) bool {
	cmd := exec.Command("java", "-cp", "scsl.jar", "TestDec.java", FileName)
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
	if fm.CheckFile("{DEC}"+FileName, "DRM Decrypt") {
		fm.FileDelete(FileName)
		return true
	} else {
		return false
	}
}

func (fm *DrmManager) DRM_CheckEnc(FileName string) int {
	cmd := exec.Command("java", "-cp", "scsl.jar", "CheckEnc.java", FileName)
	out, err := cmd.Output()
	if err != nil {
		log.Println("DRMCheck_Err", err)
		return 0
	}

	if bytes.Equal(out, []byte{49, 13, 10}) || bytes.Equal(out, []byte{49, 10}) { // 이미 암호화됐을때 구분
		log.Println("File is Encrypted")
		return 1
	} else {
		log.Println("File is not Encrypted : ", FileName, "->", string(out), "/", out)
		return 2
	}
}

func (fm *DrmManager) CheckFile(FileName string, logMessage string) bool {
	for i := 0; i < 20; i++ {
		if _, err := os.Stat(FileName); err != nil {
			log.Println(logMessage, " Fail : ", FileName, "Err -", err)
		} else {
			log.Println(logMessage, " Success : ", FileName)
			return true
		}
		time.Sleep(time.Millisecond * 500)
	}
	return false
}

func (fm *DrmManager) FileDelete(FileName string) {
	cmd := exec.Command("rm", "-rf", FileName)
	cmd.Run()
}

func (fm *DrmManager) FileGetFail(conn string, FileName string, ErrorType string) {
	// errPacket := R_Error{Status: 2}
	// instance_gs.SendPacketByConn(c, errPacket, Error)
	log.Println("FileGetFail(", ErrorType, ")-", FileName)
}

func (fm *DrmManager) PptToPdf(FileName string, FolderPath string) bool {
	cmd := exec.Command("unoconv", "-f", "pdf", FileName)
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}

	if fm.CheckFile(FileName[:len(FileName)-4]+"pdf", "PPT Convert") {
		fm.FileDelete(FileName)
		return true
	} else {
		return false
	}
}

func (fm *DrmManager) FileChangeNameMove(FileName string, FilePath string, NasPath string) bool {
	cmd := exec.Command("mv", "-f", FilePath+FileName, NasPath)
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}

	return fm.CheckFile(NasPath, "FileChangeNameMove")
}
