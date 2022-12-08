package main

import (
	"log"
	"net/url"
	"strings"

	cts "github.com/DW-inc/FileServer/CTS"
	drm "github.com/DW-inc/FileServer/DRM"
	logm "github.com/DW-inc/FileServer/Log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var Port, NasPath, CTSAddress string

func main() {
	ServerTypeSetting(0) // 0: PPRK, 1: HyunDai STR 8088, 2: HyunDai CTS 4401
	//------------ INIT Setting  ------------//
	logm.GetLogManager().SetLogFile(Port)
	app := fiber.New(fiber.Config{
		BodyLimit: 9999 * 1024 * 1024,
		// ie
		// JSONEncoder: json.Marshal,
		// JSONDecoder: json.Unmarshal,
	})
	app.Use(cors.New(cors.ConfigDefault))
	app.Use(logger.New(logger.ConfigDefault))
	cts.GetCtsManager().Init(CTSAddress)
	//CtsConn := ConnecToCTS()
	//------------ INIT Setting  ------------//

	app.Static("/uploadpage", "./UploadPage")

	// app.Use("/nas", filesystem.New(filesystem.Config{
	// 	Root: http.Dir(NasPath),
	// }))

	app.Get("/nas/:ChNum/:FileName", func(c *fiber.Ctx) error {
		//c.Set(fiber.HeaderCacheControl, "no-store, max-age=0")
		//log.Println(c.Request().Header.String())
		//log.Println(c.Response().Header.String())

		//target := c.Get("Target")
		//log.Println(target)
		log.Println(c.Params("ChNum") + "/" + c.Params("FileName"))
		filePath := c.Params("ChNum") + "/" + c.Params("FileName")

		filePath, err := url.QueryUnescape(filePath)
		if err != nil {
			log.Println("url parse faile :", err)
		}
		log.Println("parsee", filePath)

		// jData, err := json.Marshal(filePath)
		// if err != nil {
		// 	log.Println(err)
		// }
		// log.Println("JSONNNN:", string(jData))

		return c.Download("../Server/Storage/nas/" + filePath)
		//return c.Download()
	})

	app.Post("/upload", func(c *fiber.Ctx) error {
		//log.Println("BODY:", string(c.Body()))
		//log.Println("HEADER1", string(c.Request().Header.Header()))
		//log.Println("HEADER2", string(c.Response().Header.Header()))

		LocalIP := strings.Split(c.Context().RemoteAddr().String(), ":")[0]

		c.Context().SetContentType("multipart/form-data")
		if file, err := c.FormFile("file"); err != nil {
			log.Println("upload fail", err)
		} else {
			if SaveFileErr := c.SaveFile(file, file.Filename); SaveFileErr != nil {
				log.Println("SaveFile fail", SaveFileErr)
			} else {
				// Send Packet to CTS Server

				TempFileName := file.Filename
				FinalFileName := file.Filename
				IsPPT := "pptx" == TempFileName[len(TempFileName)-4:]
				conn := "TestID"
				// 파일 암호화여부 체크
				IsEnc := drm.GetDrmManager().DRM_CheckEnc(TempFileName)
				switch IsEnc { // 0: CheckFail, 1: Encrypted, 2: NotEncrypted
				case 0:
					drm.GetDrmManager().FileGetFail(conn, TempFileName, "DRM_CheckFile")
					return nil
				case 1: //암호화 되었을때
					if drm.GetDrmManager().DRM_Decrypt(TempFileName) { // 복호화
						TempFileName = "{DEC}" + TempFileName
						if IsPPT {
							if drm.GetDrmManager().PptToPdf(TempFileName, "") { // PPT였다면 PDF변환
								TempFileName = TempFileName[:len(TempFileName)-4] + "pdf"
								FinalFileName = FinalFileName[:len(FinalFileName)-4] + "pdf"
								IsEnc = 2
							} else {
								drm.GetDrmManager().FileGetFail(conn, TempFileName, "PDF_Convert")
								return nil
							}
						}
					} else {
						drm.GetDrmManager().FileGetFail(conn, TempFileName, "DRM_Decrypt")
						return nil
					}
				case 2: // 암호화 안되어 있을때
					if IsPPT {
						if drm.GetDrmManager().PptToPdf(TempFileName, "") {
							TempFileName = TempFileName[:len(TempFileName)-4] + "pdf"
							FinalFileName = FinalFileName[:len(FinalFileName)-4] + "pdf"
						} else {
							// PDF변환 실패 알리기
							drm.GetDrmManager().FileGetFail(conn, TempFileName, "PDF_Convert")

							// PDF변환 실패해도 progress바 없애줘야해서 줌
							// packet2 := R_UploadComplete{}
							// packet2.IsFinished = true
							// sendBuffer := MakeSendBuffer(EUploadComplete, packet2)
							// instance_gs.SendByte(conn, sendBuffer)
							return nil
						}
					}
				}

				// NasUpload
				Success := drm.GetDrmManager().FileChangeNameMove(TempFileName, "", NasPath+"{"+LocalIP+"}"+FinalFileName)

				// 업로드 성공 패킷 Send
				// packet2 := R_UploadComplete{}
				// packet2.IsFinished = true
				// sendBuffer := MakeSendBuffer(EUploadComplete, packet2)
				// instance_gs.SendByte(conn, sendBuffer)
				if Success {
					// FileList Update
					// db := &structs.FileList{Channel: ChannelNum, FileName: FileName, Size: int32(len(*f))}
					// fl := &structs.FileList{}
					// if r := GetContentManager().DBMS.Table("file_list").Where("channel = ? AND file_name = ?", ChannelNum, FileName).Find(fl); r.RowsAffected == 0 {
					// 	GetContentManager().DBMS.Create(db)
					// }
					// instance_gs.SendPacketByConn(conn, GetContentManager().GetFileList(ChannelNum, false), FileList)

					//time.Sleep(time.Second * 1)

					// Send PPT 변환파일
					// if IsPPT {
					// 	log.Println("ppt변환 끝 보내기 시작")
					// 	signalPacket := R_DownloadPPTtoPDF{IsStart: true, FileName: FileName}
					// 	instance_gs.SendPacketByConn(conn, signalPacket, EDownloadPPTtoPDF)
					// 	if fm.SendFile(conn, FileName, false) {
					// 		log.Println("ppt변환 끝 보내기 성공")
					// 	} else {
					// 		log.Println("ppt변환 끝 보내기 실패")
					// 	}
					// }
					LocalIP := strings.Split(c.Context().RemoteAddr().String(), ":")[0]

					packet := S_WebFileCompelete{Ip: LocalIP, IsSuccess: true, FileName: FinalFileName}
					cts.GetCtsManager().SendPacket(packet, 89)
					log.Println("File Get Success : ", FinalFileName)

				} else {
					// errPacket := R_Error{Status: 2}
					// instance_gs.SendPacketByConn(conn, errPacket, Error)

					packet := S_WebFileCompelete{Ip: LocalIP, IsSuccess: false, FileName: FinalFileName}
					cts.GetCtsManager().SendPacket(packet, 89)
					log.Println("File Get Fail: ", file.Filename)
				}
			}
		}
		return nil
	})

	app.Listen(Port)
}

type S_WebFileCompelete struct {
	Ip        string
	IsSuccess bool
	FileName  string
}

func ServerTypeSetting(Stype int) { // 0: PPRK, 1: HyunDai STR 8088, 2: HyunDai CTS 4401
	switch Stype {
	case 0:
		Port = ":8009"
		NasPath = "../Server/Storage/nas/"
		CTSAddress = "192.168.0.9:8001"
	case 1:
		Port = ":8088"
		NasPath = "/dipnas/DIPServer/Storage/"
		CTSAddress = "10.5.147.88:8000"
	case 2:
		Port = ":4401"
		NasPath = "/dipnas/DIPServer/Storage/"
		CTSAddress = "10.5.147.88:8000"
	}
}
