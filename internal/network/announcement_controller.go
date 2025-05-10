package network

import (
	"context"
	"net"

	"github.com/bigcubecat/netsnake/internal/config"
	protocol "github.com/bigcubecat/netsnake/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type AnnouncementController struct {
	ctx           *context.Context
	conn          net.Conn // соединение для прослушивания
	address       net.Addr // мультикаст адрес
	MulticastAddr string

	inboxMessageQueue chan MessagePromise
}

func NewAnnouncementController(ctx *context.Context, address string) *AnnouncementController {
	return &AnnouncementController{
		ctx:               ctx,
		MulticastAddr:     address,
		inboxMessageQueue: make(chan MessagePromise, 100),
	}
}

func (ac *AnnouncementController) InboxRoutine() {
	logrus.Println("InboxRoutine")
	go func() {
		for {
			select {
			case <-(*ac.ctx).Done():
				logrus.Println("multicast listener finished")
				return
			default:
				ac.process()
				return
			}
		}
	}()
}

func (ac *AnnouncementController) process() {
	addr, err := net.ResolveUDPAddr("udp", ac.MulticastAddr)
	if err != nil {
		logrus.Fatalf("Ошибка при разрешении адреса: %v", err)
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		logrus.Fatalf("Ошибка при прослушивании: %v", err)
	}
	defer conn.Close()

	buf := make([]byte, config.BUFFER_SIZE)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
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
		ac.inboxMessageQueue <- MessagePromise{Message: &msg, Address: remoteAddr.IP.String(), Port: remoteAddr.Port}
	}
}

func (ac *AnnouncementController) ReadInbox() []MessagePromise {
	n := len(ac.inboxMessageQueue)
	result := make([]MessagePromise, n)
	for i := 0; i < n; i++ {
		result[i] = <-ac.inboxMessageQueue
	}
	return result
}
