package network

import (
	"net"
	"time"

	"github.com/bigcubecat/netsnake/internal/master"
	"github.com/bigcubecat/netsnake/internal/network/message"
	snakes "github.com/bigcubecat/netsnake/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

const (
	multicastAddr = "239.192.0.4:9192"
)

type Network struct {
	Conn   *net.UDPConn
	Master *master.Master
	MsgSeq int64
}

func NewNetwork(config *snakes.GameConfig) *Network {
	return &Network{
		Master: master.NewMaster(config),
	}
}

func (nw *Network) Start() error {
	addr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		return err
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	nw.Conn = conn

	go nw.listen()
	go nw.sendAnnouncements()
	return nil
}

func (nw *Network) listen() {
	buf := make([]byte, 1024)
	for {
		nw.Conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, addr, err := nw.Conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		msg, err := message.DecodeGameMessage(buf[:n])
		if err != nil {
			continue
		}

		nw.HandleMessage(msg, addr)
	}
}

func (n *Network) HandleMessage(msg *snakes.GameMessage, addr *net.UDPAddr) {
	switch v := msg.Type.(type) {
	case *snakes.GameMessage_Join:
		logrus.Debug("message type=", v)
		n.HandleJoinRequest(msg, addr) // Передаем полное сообщение
		// Обработка других типов сообщений
	}
}

func (n *Network) HandleJoinRequest(msg *snakes.GameMessage, addr *net.UDPAddr) {
	join := msg.GetJoin() // Получаем JoinMsg из сообщения
	if join == nil {
		return
	}

	playerID, err := n.Master.HandleJoin(join, addr)
	if err != nil {
		n.sendError(msg, addr, err.Error())
		return
	}

	n.sendAck(addr, playerID)
}

func (n *Network) sendAck(addr *net.UDPAddr, playerID int32) {
	ack := &snakes.GameMessage{
		MsgSeq:     proto.Int64(n.MsgSeq),
		SenderId:   proto.Int32(0), // ID MASTER
		ReceiverId: proto.Int32(playerID),
		Type: &snakes.GameMessage_Ack{
			Ack: &snakes.GameMessage_AckMsg{},
		},
	}

	data, _ := message.EncodeGameMessage(ack)
	n.Conn.WriteToUDP(data, addr)
	n.MsgSeq++
}

func (n *Network) sendError(
	msg *snakes.GameMessage,
	addr *net.UDPAddr,
	errorMsg string,
) {
	errMsg := &snakes.GameMessage{
		MsgSeq:     proto.Int64(n.MsgSeq),
		SenderId:   proto.Int32(0), // ID MASTER
		ReceiverId: msg.ReceiverId, // Используем ReceiverId из оригинального сообщения
		Type: &snakes.GameMessage_Error{
			Error: &snakes.GameMessage_ErrorMsg{
				ErrorMessage: proto.String(errorMsg),
			},
		},
	}

	data, _ := message.EncodeGameMessage(errMsg) // Кодируем GameMessage
	n.Conn.WriteToUDP(data, addr)
	n.MsgSeq++
}

func (n *Network) sendAnnouncements() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Формирование и отправка анонсов
	}
}
