package network

import (
	"context"

	"github.com/bigcubecat/netsnake/internal/common"
	"github.com/bigcubecat/netsnake/internal/config"
	"github.com/bigcubecat/netsnake/internal/model"
	"github.com/bigcubecat/netsnake/internal/utils"
	protocol "github.com/bigcubecat/netsnake/proto"
	"github.com/sirupsen/logrus"
)

type Peer struct {
	ID      int
	Role    protocol.NodeRole
	Players map[int]common.Player

	GameInstance *model.Game
	Config       *config.Config

	annoncementController AnnouncementController
	messageController     MessageController

	ctx    context.Context
	cancel context.CancelFunc
}

func NewPeer(g *model.Game, conf *config.Config) *Peer {
	peer := &Peer{
		ID:           utils.RandomId(),
		Role:         modeToRole(conf.UiConfig.Mode),
		Players:      make(map[int]common.Player),
		GameInstance: g,
		Config:       conf,
	}
	peer.Players[peer.ID] = common.Player{
		ID:    peer.ID,
		Name:  peer.Config.CliConfig.PlayerName,
		Score: 0,
	}
	peer.GameInstance.AddSnake(peer.ID, model.Master)
	peer.annoncementController = *NewAnnouncementController(
		&peer.ctx,
		conf.CliConfig.MulticastAddress,
	)
	peer.messageController = *NewMessageController(
		&peer.ctx,
		peer,
	)
	return peer
}

func (peer *Peer) Exit() {
	peer.cancel()
}

func (peer *Peer) StartGorutines() {
	peer.ctx, peer.cancel = context.WithCancel(context.Background())

	go peer.routine()
	go peer.annoncementController.InboxRoutine()
	go peer.messageController.Routine()
}

func (peer *Peer) routine() {
	for {
		logrus.Println("routine")
		select {
		case <-peer.ctx.Done():
			logrus.Println("peer routine done")
			return
		default:
			// таймаут в Process так как у разных ролей он разный
			GetStateByRole(peer.Role).Process(peer)
		}
	}
}

func modeToRole(mode int) protocol.NodeRole {
	switch mode {
	case 0:
		return protocol.NodeRole_MASTER
	case 1:
		return protocol.NodeRole_NORMAL
	default:
		return protocol.NodeRole_VIEWER
	}
}
