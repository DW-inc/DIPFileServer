package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	logm "github.com/DW-inc/FileServer/Log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

func main() {
	//------------ INIT Setting  ------------//
	logm.GetLogManager().SetLogFile()
	app := fiber.New(fiber.Config{
		BodyLimit: 32 * 1024 * 102,
	})
	app.Use(cors.New(cors.ConfigDefault))

	ConnecToCTS()
	//------------ INIT Setting  ------------//

	app.Static("/uploadpage", "./UploadPage")

	app.Use("/nas", filesystem.New(filesystem.Config{
		Root: http.Dir("../Server/Storage/nas"),
	}))

	app.Post("/upload", func(c *fiber.Ctx) error {
		log.Println("BODY:", string(c.Body()))
		//log.Println("HEADER1", string(c.Request().Header.Header()))
		//log.Println("HEADER2", string(c.Response().Header.Header()))
		c.Context().SetContentType("multipart/form-data")
		if file, err := c.FormFile("file"); err != nil {
			log.Println("upload fail", err)
		} else {
			if SaveFileErr := c.SaveFile(file, fmt.Sprint("../Server/Storage/", file.Filename)); SaveFileErr != nil {
				log.Println("SaveFile fail", SaveFileErr)
			} else {
				// Send Packet to CTS Server

			}
		}
		return nil
	})

	app.Listen(":8009")
}

type S_WebFileCompelete struct {
	Id        string
	IsSuccess bool
}

func ConnecToCTS() {
	//conn, err := net.Dial("tcp", "192.168.0.9:8001")
	//if err != nil {
	//	fmt.Println("Faield to Dial : ", err)
	//}
	////defer conn.Close()
	//
	//go func(c net.Conn) {
	//	packet := S_WebFileCompelete{Id: "tester", IsSuccess: true}
	//	sendBuffer := MakeSendBuffer(89, packet)
	//
	//	_, err = c.Write(sendBuffer)
	//	//test := "etetsteesettes"
	//	//_, err = c.Write([]byte(test))
	//	if err != nil {
	//		fmt.Println("Failed to write data : ", err)
	//	}
	//
	//}(conn)

}

func MakeSendBuffer[T any](pktid uint16, data T) []byte {
	sendData, err := json.Marshal(&data)
	if err != nil {
		log.Println("MakeSendBuffer : Marshal Error", err)
	}
	sendBuffer := make([]byte, 4)

	pktsize := len(sendData) + 4

	binary.LittleEndian.PutUint16(sendBuffer, uint16(pktsize))
	binary.LittleEndian.PutUint16(sendBuffer[2:], pktid)

	sendBuffer = append(sendBuffer, sendData...)

	return sendBuffer
}
