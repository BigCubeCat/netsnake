package network

import (
	"github.com/bigcubecat/netsnake/internal/model"
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
		state.sendStateMsg(peer)
	}
	state.sendAnnMsg(peer) // делаем спам рассылку с новостями

	multicastInbox := peer.annoncementController.ReadInbox()
	for _, recvMessage := range multicastInbox {
		switch recvMessage.Message.GetType().(type) {
		case *protocol.GameMessage_Discover:
			// TODO: сделать операцию на определение DEPUTY
			logrus.Debug("discover message recv")
			// делаем спам рассылку с новостями
			state.sendAnnMsgUnicast(peer, recvMessage.Address, recvMessage.Port)
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
	logrus.Debugln(
		"state master process ",
		multicastInbox,
		unicastInbox,
	)
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

// рассылка сообщения по мультикасту
func (state PeerMasterState) sendAnnMsgUnicast(peer *Peer, address string, port int) {
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
	peer.messageController.AddMessage(address, port, msg)
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

func (state PeerMasterState) sendStateMsg(peer *Peer) {
	logrus.Debug("send state msg")
	peer.stateOrder++
	for id, player := range peer.Players {
		if player.Port == 0 {
			continue
		}
		msg := message.NewStateMsg(
			peer.msgSeq.Load(),
			int32(peer.ID),
			int32(id),
			peer.stateOrder,
			state.generateSnakes(peer),
			state.generateFood(peer),
			state.generatePlayers(peer),
		)

		logrus.Debugln(id, "send message to ", player.IpAddress, player.Port, player.Role, msg)
		peer.messageController.AddMessage(
			peer.Config.CliConfig.MulticastAddress,
			peer.Config.CliConfig.MulticastPort,
			msg,
		)
	}
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

func (state PeerMasterState) generateFood(peer *Peer) []*protocol.GameState_Coord {
	foods := make([]*protocol.GameState_Coord, 0)
	for _, point := range peer.GameInstance.State.Food {
		foods = append(foods, &protocol.GameState_Coord{
			X: proto.Int32(int32(point.X)),
			Y: proto.Int32(int32(point.Y)),
		})
	}
	return foods
}

func (state PeerMasterState) generateSnakes(peer *Peer) []*protocol.GameState_Snake {
	snakes := make([]*protocol.GameState_Snake, 0)
	for _, snake := range peer.GameInstance.State.Snakes {
		body := make([]*protocol.GameState_Coord, len(snake.Body))
		for i := 0; i < len(body); i++ {
			body[i] = &protocol.GameState_Coord{
				X: proto.Int32(int32(snake.Body[i].X)),
				Y: proto.Int32(int32(snake.Body[i].Y)),
			}
		}
		snakes = append(snakes, &protocol.GameState_Snake{
			PlayerId:      proto.Int32(int32(snake.ID)),
			HeadDirection: snakeDir(snake.Direction),
			Points:        body,
		})
	}
	return snakes
}

func snakeDir(dir model.Direction) *protocol.Direction {
	switch dir {
	case model.Up:
		return protocol.Direction_UP.Enum()
	case model.Down:
		return protocol.Direction_DOWN.Enum()
	case model.Right:
		return protocol.Direction_RIGHT.Enum()
	default:
		return protocol.Direction_LEFT.Enum()
	}
}
