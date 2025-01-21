package network

import (
	"time"

	"github.com/bigcubecat/netsnake/internal/network/message"
	protocol "github.com/bigcubecat/netsnake/proto"
	"github.com/sirupsen/logrus"
)

func (state PeerMasterState) Process(peer *Peer) {
	if peer.GameInstance == nil {
		logrus.Println("game inst is null")
		logrus.Fatalln("game instance is nill in master")
	}
	peer.GameInstance.MoveSnakes()
	state.sendAnnMsg(peer) // делаем спам рассылку с новостями

	multicastInbox := peer.annoncementController.ReadInbox()
	logrus.Println("messages recieaved: ", len(multicastInbox))
	for _, recvMessage := range multicastInbox {
		switch recvMessage.GetType().(type) {
		case *protocol.GameMessage_Discover:
			// TODO: сделать операцию на определение DEPUTY
			logrus.Debug("discover message recv")
			state.sendAnnMsg(peer) // делаем спам рассылку с новостями
		}
	}
	time.Sleep(time.Duration(peer.Config.EnvConfig.Dt) * time.Millisecond)
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
	peer.messageController.AddMessage(peer.Config.CliConfig.MulticastAddress, msg)
}

func (state PeerMasterState) generatePlayers(peer *Peer) *protocol.GamePlayers {
	return &protocol.GamePlayers{Players: []*protocol.GamePlayer{}}
}
