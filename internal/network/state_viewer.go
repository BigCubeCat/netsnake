package network

import (
	"fmt"

	"github.com/bigcubecat/netsnake/internal/common"
	"github.com/bigcubecat/netsnake/internal/model"
	protocol "github.com/bigcubecat/netsnake/proto"
	"github.com/sirupsen/logrus"
)

func (state PeerViewerState) Process(peer *Peer) {
	multicastInbox := peer.annoncementController.ReadInbox()
	unicastInbox := peer.messageController.ReadInbox()
	for _, recvMessage := range multicastInbox {
		switch recvMessage.Message.GetType().(type) {
		case *protocol.GameMessage_Announcement:
			state.updateGame(peer, recvMessage.Message.GetAnnouncement())
		case *protocol.GameMessage_State:
			fmt.Println("STATE")
			state.updateState(peer, recvMessage.Message.GetState())
		}
	}
	for _, recvMessage := range unicastInbox {
		switch recvMessage.Message.GetType().(type) {
		case *protocol.GameMessage_Announcement:
			state.updateGame(peer, recvMessage.Message.GetAnnouncement())
		case *protocol.GameMessage_State:
			logrus.Println("STATE UPDATED")
			state.updateState(peer, recvMessage.Message.GetState())
		}
	}
	logrus.Debugln(
		"state view process ",
		multicastInbox,
		unicastInbox,
	)
}

func (state PeerViewerState) updateState(peer *Peer, data *protocol.GameMessage_StateMsg) {
	players := data.GetState().Players.Players
	for k := range peer.Players {
		delete(peer.Players, k)
	}
	for _, player := range players {
		id := int(player.GetId())
		peer.Players[id] = common.Player{
			ID:        id,
			Role:      roleToMode(*player.Role.Enum()),
			Name:      player.GetName(),
			IpAddress: player.GetIpAddress(),
			Port:      int(player.GetPort()),
			Score:     int(player.GetScore()),
		}
	}

	peer.GameInstance.ResetSnakes()
	snakes := data.GetState().GetSnakes()
	for _, snake := range snakes {
		points := snake.GetPoints()
		body := make([]model.Point, len(points))
		for i := 0; i < len(points); i++ {
			body[i] = model.Point{X: int(points[i].GetX()), Y: int(points[i].GetY())}
		}
		dir := snake.GetHeadDirection()
		isAlive := modeToRole(
			peer.Players[int(snake.GetPlayerId())].Role,
		) != protocol.NodeRole_VIEWER
		newSnake := model.Snake{
			Direction: model.Direction(dir),
			Body:      body,
			Score:     peer.Players[int(snake.GetPlayerId())].Score,
			IsAlive:   isAlive,
		}
		peer.GameInstance.SetSnake(&newSnake)
	}
}

func (state PeerViewerState) updateGame(peer *Peer, data *protocol.GameMessage_AnnouncementMsg) {
	if peer.GameInstance == nil {
		logrus.Fatal("GameInstance is nil")
	}
	games := data.GetGames()
	if len(games) == 0 {
		return
	}
	game := games[0]
	peer.GameInstance.Rebuild(
		int(game.Config.GetWidth()),
		int(game.Config.GetHeight()),
		int(game.Config.GetFoodStatic()),
	)
}
