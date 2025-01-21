package network

import (
	"context"
	"fmt"
	"net"

	"github.com/bigcubecat/netsnake/internal/config"
	"github.com/bigcubecat/netsnake/internal/network/message"
	protocol "github.com/bigcubecat/netsnake/proto"
	"github.com/sirupsen/logrus"
)

type MessageController struct {
	conn    *net.UDPConn
	ctx     *context.Context
	peerPtr *Peer

	inboxMessageQueue chan MessagePromise

	outboxMessageQueue chan MessagePromise
}

func NewMessageController(
	ctx *context.Context,
	conn *net.UDPConn,
	peerPtr *Peer,
) *MessageController {
	return &MessageController{
		conn:    conn,
		ctx:     ctx,
		peerPtr: peerPtr,

		inboxMessageQueue:  make(chan MessagePromise, 1000),
		outboxMessageQueue: make(chan MessagePromise, 1000),
	}
}

func (mc *MessageController) Routine() {
	logrus.Println("MessageController Routine")
	go mc.recvData()
	go mc.sendData()
}

// планируем отправку сообщения при первой же возможности
func (mc *MessageController) AddMessage(
	address string,
	port int,
	gameMessage *protocol.GameMessage,
) {
	mc.peerPtr.msgSeq.Add(1)
	mc.outboxMessageQueue <- MessagePromise{
		Address: address,
		Port:    port,
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
					&net.UDPAddr{IP: net.ParseIP(msg.Address), Port: msg.Port},
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
			msg, err := message.UnmarshalGameMessage(buffer[:n])
			if err != nil {
				logrus.Errorf("cant UnmarshalGameMessage %s: %s\n",
					remoteAddr,
					string(buffer[:n]),
				)
				continue
			}
			fmt.Println("получено сообщение")
			mc.inboxMessageQueue <- MessagePromise{
				Address: remoteAddr.IP.String(),
				Port:    remoteAddr.Port,
				Message: msg,
			}
		}
	}
}

func (mc *MessageController) ReadInbox() []MessagePromise {
	n := len(mc.inboxMessageQueue)
	result := make([]MessagePromise, n)
	for i := 0; i < n; i++ {
		result[i] = <-mc.inboxMessageQueue
	}
	return result
}
