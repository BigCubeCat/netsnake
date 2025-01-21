package network

import (
	// "github.com/bigcubecat/netsnake/internal/model"
	"github.com/bigcubecat/netsnake/internal/network/message"
	protocol "github.com/bigcubecat/netsnake/proto"
	"github.com/sirupsen/logrus"
)

func (state PeerJoinState) Process(peer *Peer) {
	msg := message.NewDiscoverMsg(peer.msgSeq.Load())
	peer.messageController.AddMessage(
		peer.Config.CliConfig.MulticastAddress,
		peer.Config.CliConfig.MulticastPort,
		msg,
	)
	multicastInbox := peer.annoncementController.ReadInbox()
	unicastInbox := peer.messageController.ReadInbox()
	logrus.Debugln(
		"state join process ",
		multicastInbox,
		unicastInbox,
	)
	for _, recvMessage := range unicastInbox {
		switch recvMessage.Message.GetType().(type) {
		case *protocol.GameMessage_Ack:
			state.joinPeer(peer)
			return
		case *protocol.GameMessage_Announcement:
			state.handleAnnouncement(peer, recvMessage)
			break
		}
	}
	for _, recvMessage := range multicastInbox {
		switch recvMessage.Message.GetType().(type) {
		case *protocol.GameMessage_Announcement:
			peer.masterFound = true
			break
		}
	}
}

func (state PeerJoinState) handleAnnouncement(
	peer *Peer,
	msg MessagePromise,
) {
	// УРА! ПОДКЛЮЧАЕМСЯ
	// к сожалению, до сдачи лабы менее 10 часов,
	// да и по протоколу мастер в сети только один,
	// так что: автоподключение
	peer.messageController.AddMessage(
		msg.Address,
		msg.Port,
		message.CreateJoinMessage(
			peer.Config.CliConfig.PlayerName,
			peer.Config.CliConfig.GameName,
			peer.Role,
		),
	)
	if len(msg.Message.GetAnnouncement().GetGames()) == 0 {
		return
	}
	// g := msg.Message.GetAnnouncement().GetGames()[0]
	// peer.GameInstance = model.NewGame(
	// 	int(g.Config.GetWidth()),
	// 	int(g.Config.GetWidth()),
	// 	int(g.Config.GetFoodStatic()),
	// )
}

func (state PeerJoinState) joinPeer(peer *Peer) {
	peer.needJoin = false
}
