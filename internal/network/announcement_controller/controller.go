package announcementcontroller

import (
	"context"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type AnnouncementController struct {
	Conn          net.Conn        // соединение для прослушивания
	Address       net.Addr        // мультикаст адрес
	Ctx           context.Context // контест горутины слушателя
	multicastAddr string
}

func (ac *AnnouncementController) InboxRoutine() {
	addr, err := net.ResolveUDPAddr("udp", ac.multicastAddr)
	if err != nil {
		logrus.Fatalf("Ошибка при разрешении адреса: %v", err)
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	conn.SetWriteBuffer(1024)
	if err != nil {
		logrus.Fatalf("Ошибка при прослушивании: %v", err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			logrus.Printf("Ошибка при чтении: %v", err)
			continue
		}

		var msg protocol.Message
		err = proto.Unmarshal(buf[:n], &msg)
		if err != nil {
			logrus.Printf("Ошибка десериализации: %v", err)
			continue
		}
	}
}
