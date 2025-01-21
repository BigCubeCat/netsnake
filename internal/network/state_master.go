package network

import (
	"github.com/bigcubecat/netsnake/internal/network/message"
	protocol "github.com/bigcubecat/netsnake/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func (state PeerMasterState) Process(peer *Peer) {
	if peer.GameInstance == nil {
		logrus.Println("game inst is null")
		logrus.Fatalln("game instance is nill in master")
	}
	if peer.step == 0 {
		peer.GameInstance.MoveSnakes()
	}
	state.sendAnnMsg(peer) // делаем спам рассылку с новостями

	multicastInbox := peer.annoncementController.ReadInbox()
	for _, recvMessage := range multicastInbox {
		switch recvMessage.GetType().(type) {
		case *protocol.GameMessage_Discover:
			// TODO: сделать операцию на определение DEPUTY
			logrus.Debug("discover message recv")
			state.sendAnnMsg(peer) // делаем спам рассылку с новостями
		}
	}
	unicastInbox := peer.messageController.ReadInbox()
	for _, recvMessage := range unicastInbox {
		switch recvMessage.Message.GetType().(type) {
		case *protocol.GameMessage_Join:
			logrus.Debug("join message recv")
			peer.JoinPlayer(
				recvMessage.Message.GetJoin(),
				recvMessage.Address,
				recvMessage.Port,
			)
		}
	}
}

// рассылка сообщения по мультикасту
func (state PeerMasterState) sendAnnMsg(peer *Peer) {
	msg := message.NewAnnouncementMessage(
		peer.Config.CliConfig.PlayerName,
		peer.Config.CliConfig.GameName,
		int32(peer.Config.EnvConfig.FieldWidth),
		int32(peer.Config.EnvConfig.FieldHeight),
		int32(peer.Config.EnvConfig.FoodStatic),
		int32(peer.Config.EnvConfig.Dt),
		peer.Role,
		state.generatePlayers(peer),
	)
	peer.messageController.AddMessage(
		peer.Config.CliConfig.MulticastAddress,
		peer.Config.CliConfig.MulticastPort,
		msg,
	)
}

// Ack msg
func (state PeerMasterState) sendAckMsg(peer *Peer, recvId int) {
	msg := message.NewAckMsg(peer.msgSeq.Load(), int32(peer.ID), int32(recvId))
	peer.messageController.AddMessage(
		peer.Config.CliConfig.MulticastAddress,
		peer.Config.CliConfig.MulticastPort,
		msg,
	)
}

func (state PeerMasterState) generatePlayers(peer *Peer) *protocol.GamePlayers {
	players := make([]*protocol.GamePlayer, 0)
	for _, player := range peer.Players {
		players = append(players, &protocol.GamePlayer{
			Name:      proto.String(player.Name),
			Id:        proto.Int32(int32(peer.ID)),
			IpAddress: proto.String(player.IpAddress),
			Port:      proto.Int32(int32(player.Port)),
			Role:      modeToRole(player.Role).Enum(),
			Type:      protocol.Default_GamePlayer_Type.Enum(),
			Score:     proto.Int32(int32(player.Score)),
		})
	}

	return &protocol.GamePlayers{Players: []*protocol.GamePlayer{}}
}
