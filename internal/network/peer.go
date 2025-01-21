package network

import (
	"context"
	"net"
	"time"

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

	step uint8 // шаг [0, 10)
}

func NewPeer(g *model.Game, conf *config.Config) *Peer {
	peer := &Peer{
		ID:           utils.RandomId(),
		Role:         modeToRole(conf.UiConfig.Mode),
		Players:      make(map[int]common.Player),
		GameInstance: g,
		Config:       conf,
		step:         0,
	}
	addr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		logrus.Fatalf("ошибка при разрешении адреса: %s", err.Error())
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		logrus.Fatalf("bind error: %s", err.Error())
	}
	peer.Players[peer.ID] = common.Player{
		ID:        peer.ID,
		Name:      peer.Config.CliConfig.PlayerName,
		IpAddress: addr.String(),
		Port:      addr.Port,
		Score:     0,
	}
	if peer.Role == protocol.NodeRole_MASTER {
		peer.GameInstance.AddSnake(peer.ID, model.Master)
	}
	peer.annoncementController = *NewAnnouncementController(&peer.ctx, conf.CliConfig.MulticastAddress)
	peer.messageController = *NewMessageController(&peer.ctx, conn, peer)
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
			peer.step = (peer.step + 1) % 10
			GetStateByRole(peer.Role).Process(peer)
			time.Sleep(time.Duration(peer.Config.EnvConfig.Dt/10) * time.Millisecond)
		}
	}
}

func (peer *Peer) SetUnicastAddress(address string, port int) {
	value := peer.Players[peer.ID]
	value.IpAddress = address
	value.Port = port
	peer.Players[peer.ID] = value
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
