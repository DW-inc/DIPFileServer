package setting

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
	"sync"
	"time"
)

type CtsHandler struct {
	IP   string
	Conn net.Conn
}

var Cts_Ins *CtsHandler
var Cts_once sync.Once

func GetCtsManager() *CtsHandler {
	Cts_once.Do(func() {
		Cts_Ins = &CtsHandler{}
	})
	return Cts_Ins
}

func (cts *CtsHandler) Init(IP string) {
	cts.IP = IP
	cts.ConnecToCTS(IP)
}

func (cts *CtsHandler) ConnecToCTS(IP string) {
	conn, err := net.Dial("tcp", IP)
	if err != nil {
		log.Println("Faield to Dial : ", err)
	} else {
		log.Println("Successed to Dial : ", conn)
	}
	cts.Conn = conn
	//defer conn.Close()
}

func MakeSendBuffer[T any](pktid uint16, data T) []byte {
	sendData, err := json.Marshal(&data)
	if err != nil {
		log.Println("MakeSendBuffer : Marshal Error", err)
	}
	sendBuffer := make([]byte, 6)

	pktsize := len(sendData) + 6

	binary.LittleEndian.PutUint32(sendBuffer, uint32(pktsize))
	binary.LittleEndian.PutUint16(sendBuffer[4:], pktid)

	sendBuffer = append(sendBuffer, sendData...)

	return sendBuffer
}

func (cts *CtsHandler) SendPacket(recvpkt any, pkttype uint16) {
	sendBuffer := MakeSendBuffer(pkttype, recvpkt)
	sent, err := cts.Conn.Write(sendBuffer)
	if err != nil || sent == 0 || cts.Conn == nil {
		log.Println("SendPacket ERROR :", err)
		// if errors.Is(err, syscall.WSAECONNRESET) {
		// }
		cts.ConnecToCTS(cts.IP)
		time.Sleep(time.Second * 3)
		cts.SendPacket(recvpkt, pkttype)
	} else {
		if sent != len(sendBuffer) {
			log.Println("[Sent diffrent size] : SENT =", sent, "BufferSize =", len(sendBuffer))
		}
	}
}
