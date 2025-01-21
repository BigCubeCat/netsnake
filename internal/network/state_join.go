package network

import (
	"fmt"

	"github.com/bigcubecat/netsnake/internal/network/message"
	protocol "github.com/bigcubecat/netsnake/proto"
)

func (state PeerJoinState) Process(peer *Peer) {
	msg := message.NewDiscoverMsg(peer.msgSeq.Load())
	peer.messageController.AddMessage(
		peer.Config.CliConfig.MulticastAddress,
		msg,
	)

	if peer.step%10 != 0 {
		return
	}
	multicastInbox := peer.annoncementController.ReadInbox()
	for _, recvMessage := range multicastInbox {
		fmt.Println("inbox")
		switch recvMessage.GetType().(type) {
		case *protocol.GameMessage_Announcement:
			// УРА! ПОДКЛЮЧАЕМСЯ
			// к сожалению, до сдачи лабы менее 10 часов,
			// да и по протоколу мастер в сети только один,
			// так что: автоподключение
			peer.messageController.AddMessage(peer.Config.CliConfig.MulticastAddress, message.CreateJoinMessage(
				peer.Config.CliConfig.PlayerName,
				peer.Config.CliConfig.GameName,
				peer.Role,
			))
			break
		case *protocol.GameMessage_Ack:
			fmt.Println("ack=", msg.MsgSeq)
			state.joinPeer(peer)
			break
		}
	}
}

func (state PeerJoinState) joinPeer(peer *Peer) {
	peer.needJoin = false
}
