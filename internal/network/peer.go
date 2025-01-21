package network

import (
	"context"
	"time"

	"github.com/bigcubecat/netsnake/internal/config"
	"github.com/bigcubecat/netsnake/internal/model"
	"github.com/bigcubecat/netsnake/internal/utils"
	protocol "github.com/bigcubecat/netsnake/proto"
)

type Peer struct {
	ID   int
	Role protocol.NodeRole

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
		GameInstance: g,
		Config:       conf,
	}
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

	peer.annoncementController.InboxRoutine()
	peer.messageController.Routine()
	go peer.routine()
}

func (peer *Peer) routine() {
	select {
	case <-peer.ctx.Done():
		return
	default:
		GetStateByRole(peer.Role).Process(peer)
		time.Sleep(time.Duration(peer.Config.EnvConfig.Dt) * time.Millisecond)
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
