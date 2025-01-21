package network

import (
	"context"
	"net"
	"sync"

	"github.com/bigcubecat/netsnake/internal/config"
	"github.com/bigcubecat/netsnake/internal/network/message"
	protocol "github.com/bigcubecat/netsnake/proto"
	"github.com/sirupsen/logrus"
)

type MessagePromise struct {
	Address string
	Message *protocol.GameMessage
}

type MessageController struct {
	conn    *net.UDPConn
	ctx     *context.Context
	peerPtr *Peer

	inboxMessage *protocol.GameMessage
	inboxMutex   sync.Mutex

	outboxMessageQueue chan MessagePromise
}

func NewMessageController(
	ctx *context.Context,
	peerPtr *Peer,
) *MessageController {
	addr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		logrus.Fatalf("ошибка при разрешении адреса: %s", err.Error())
	}
	// Создание UDP-соединения
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		logrus.Fatalf("bind error: %s", err.Error())
	}
	return &MessageController{
		conn:    conn,
		ctx:     ctx,
		peerPtr: peerPtr,
	}
}

func (mc *MessageController) Routine() {
	go mc.recvData()
	go mc.sendData()
}

// планируем отправку сообщения при первой же возможности
func (mc *MessageController) AddMessage(address string, gameMessage *protocol.GameMessage) {
	mc.outboxMessageQueue <- MessagePromise{
		Address: address,
		Message: gameMessage,
	}
}

func (mc *MessageController) sendData() {
	var err error
	buffer := make([]byte, config.BUFFER_SIZE)
	for {
		select {
		case <-(*mc.ctx).Done():
			return
		default:
			if len(mc.outboxMessageQueue) > 0 {
				msg := <-mc.outboxMessageQueue
				buffer, err = message.MarshalGameMessage(msg.Message)
				if err != nil {
					logrus.Errorf("cant MarshalGameMessage")
					continue
				}
				_, err = mc.conn.WriteToUDP(
					buffer,
					&net.UDPAddr{IP: net.ParseIP(msg.Address)},
				)

			}
		}
	}
}

// recvData получает данные из unicast соединения
// и обрабатывает сообщение в переменную
func (mc *MessageController) recvData() {
	buffer := make([]byte, config.BUFFER_SIZE)
	for {
		select {
		case <-(*mc.ctx).Done():
			return
		default:
			n, remoteAddr, err := mc.conn.ReadFromUDP(buffer)
			if err != nil {
				logrus.Println("read error:", err)
				continue
			}
			logrus.Printf("message recieved %s: %s\n",
				remoteAddr,
				string(buffer[:n]),
			)
			mc.inboxMutex.Lock()
			mc.inboxMessage, err = message.UnmarshalGameMessage(buffer[:n])
			if err != nil {
				logrus.Errorf("cant UnmarshalGameMessage %s: %s\n",
					remoteAddr,
					string(buffer[:n]),
				)
				mc.inboxMessage = nil
			}
			mc.inboxMutex.Unlock()
		}
	}
}
