package peer

import (
	"net"
	"sync"

	"github.com/bigcubecat/netsnake/internal/config"
	"github.com/bigcubecat/netsnake/internal/model"
	snakes "github.com/bigcubecat/netsnake/proto"
)

type Peer struct {
	ID           int
	Addr         *net.UDPAddr
	Role         snakes.NodeRole
	mu           sync.Mutex
	GameInstance *model.Game
	Config       *config.Config
}

func NewPeer(g *model.Game, conf *config.Config) *Peer {
	return &Peer{
		GameInstance: g,
		Config:       conf,
	}
}
