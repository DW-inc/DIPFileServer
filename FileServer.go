package main

import (
	"log"
	"net/url"
	"strconv"
	"strings"

	cts "github.com/DW-inc/FileServer/CTS"
	db "github.com/DW-inc/FileServer/DB"
	drm "github.com/DW-inc/FileServer/DRM"
	logm "github.com/DW-inc/FileServer/Log"
	"github.com/DW-inc/FileServer/setting"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	setting.GetStManager().Init()
	db.GetDBManager().Init()

	//------------ INIT Setting  ------------//
	logm.GetLogManager().SetLogFile()
	app := fiber.New(fiber.Config{
		BodyLimit: 9999 * 1024 * 1024,
		// ie
		// JSONEncoder: json.Marshal,
		// JSONDecoder: json.Unmarshal,
	})
	app.Use(cors.New(cors.ConfigDefault))
	app.Use(logger.New(logger.ConfigDefault))
	cts.GetCtsManager().Init(setting.St_Ins.CTSAddress)
	//------------ INIT Setting  ------------//
	app.Get("/uploadpage", func(c *fiber.Ctx) error {
		IP := strings.Split(c.Context().RemoteAddr().String(), ":")[0]
		log.Println("UploadPage Request", IP)

		WebLoginData := &db.WebLogin{}
		err := db.GetDBManager().DBMS.Table("web_login").Where("ip = ?", IP).Select("sso_id", "is_online").First(&WebLoginData).Error
		if err != nil || !WebLoginData.IsOnline || WebLoginData.SsoId != "" {
			if setting.GetStManager().ServerType != 0 { //회사 내부에서는 테스트 가능하게 하기 위해
				log.Println(err)
				log.Println("IsOfflineUser Give NullPage")
				return nil
			}
		}
		log.Println("IsOnlineUser Give UploadPage")
		return c.Next()
	})

	app.Static("/uploadpage", "./UploadPage")

	// app.Use("/nas", filesystem.New(filesystem.Config{
	// 	Root: http.Dir(NasPath),
	// }))

	app.Get("/nas/:ChNum/:FileName", func(c *fiber.Ctx) error {
		ChNum, _ := strconv.Atoi(c.Params("ChNum"))
		if ChNum < -1 {
			ChNum = -1 * ChNum / 10000
		}
		ChNumstring := strconv.FormatInt(int64(ChNum), 10)
		filePath := ChNumstring + "/" + c.Params("FileName")

		filePath, err := url.QueryUnescape(filePath)
		if err != nil {
			log.Println("url parse faile :", err)
		}
		log.Println(ChNumstring + "/" + filePath)

		return c.Download(setting.St_Ins.NasPath + filePath)
	})

	app.Post("/upload", func(c *fiber.Ctx) error {
		LocalIP := strings.Split(c.Context().RemoteAddr().String(), ":")[0]

		c.Context().SetContentType("multipart/form-data")
		if file, err := c.FormFile("file"); err != nil {
			log.Println("upload fail", err)
		} else {
			if SaveFileErr := c.SaveFile(file, file.Filename); SaveFileErr != nil {
				log.Println("SaveFile fail", SaveFileErr)
			} else {
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
				Success := drm.GetDrmManager().FileChangeNameMove(TempFileName, "", setting.GetStManager().NasPath+"{"+LocalIP+"}"+FinalFileName)

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

	// app.Get("/saveCookie", func(c *fiber.Ctx) error {
	// 	IP := strings.Split(c.Context().RemoteAddr().String(), ":")[0]
	// 	GameSessionStorage.Store(IP, true)
	// 	log.Println("Online User :", IP)
	// 	// sess, err := SessionStorage.Get(c)
	// 	// if err != nil {
	// 	// 	log.Println("SessionGet Err", err)
	// 	// }
	// 	// sess.Set("ssoid", sso_id)
	// 	// if err := sess.Save(); err != nil {
	// 	// 	log.Println("SessionSave Err", err)
	// 	// } else {
	// 	// 	log.Println("SessionSave Success", sso_id)
	// 	// }
	// 	return nil
	// })

	app.Listen(setting.GetStManager().Port)
}

type S_WebFileCompelete struct {
	Ip        string
	IsSuccess bool
	FileName  string
}
