package network

import (
	"context"

	"github.com/bigcubecat/netsnake/internal/config"
	"github.com/bigcubecat/netsnake/internal/model"
	ac "github.com/bigcubecat/netsnake/internal/network/announcement_controller"
	"github.com/bigcubecat/netsnake/internal/utils"
	protocol "github.com/bigcubecat/netsnake/proto"
	"github.com/sirupsen/logrus"
)

type Peer struct {
	ID   int
	Role protocol.NodeRole

	GameInstance *model.Game
	Config       *config.Config

	annoncementController ac.AnnouncementController
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

		annoncementController: *ac.NewAnnouncementController(
			conf.CliConfig.MulticastAddress,
		),
	}
	peer.messageController = *NewMessageController(&peer.ctx, peer)
	return peer
}

func (peer *Peer) Exit() {
	peer.cancel()
}

func (peer *Peer) StartGorutines() {
	peer.ctx, peer.cancel = context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-peer.ctx.Done():
				logrus.Println("multicast listener finished")
				return
			default:
				peer.annoncementController.InboxRoutine()
				return
			}
		}
	}()
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
