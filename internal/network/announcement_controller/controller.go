package announcementcontroller

import (
	"net"
	"sync"

	"github.com/bigcubecat/netsnake/internal/config"
	protocol "github.com/bigcubecat/netsnake/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type AnnouncementController struct {
	conn          net.Conn // соединение для прослушивания
	address       net.Addr // мультикаст адрес
	MulticastAddr string
	lock          sync.Mutex

	messages []*protocol.GameMessage
}

func NewAnnouncementController(address string) *AnnouncementController {
	return &AnnouncementController{
		MulticastAddr: address,
	}
}

func (ac *AnnouncementController) InboxRoutine() {
	addr, err := net.ResolveUDPAddr("udp", ac.MulticastAddr)
	if err != nil {
		logrus.Fatalf("Ошибка при разрешении адреса: %v", err)
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	conn.SetWriteBuffer(config.BUFFER_SIZE)
	if err != nil {
		logrus.Fatalf("Ошибка при прослушивании: %v", err)
	}
	defer conn.Close()

	buf := make([]byte, config.BUFFER_SIZE)
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			logrus.Printf("Ошибка при чтении: %v", err)
			continue
		}

		var msg protocol.GameMessage
		err = proto.Unmarshal(buf[:n], &msg)
		if err != nil {
			logrus.Printf("Ошибка десериализации: %v", err)
			continue
		}
		ac.lock.Lock()
		ac.messages = append(ac.messages, &msg)
		ac.lock.Unlock()
	}
}

func (ac *AnnouncementController) ReadInbox() []*protocol.GameMessage {
	ac.lock.Lock()
	defer ac.lock.Unlock()
	result := ac.messages
	ac.messages = []*protocol.GameMessage{}
	return result
}
