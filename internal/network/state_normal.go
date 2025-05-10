package network

import (
	protocol "github.com/bigcubecat/netsnake/proto"
	"github.com/sirupsen/logrus"
)

func (state PeerNormalState) Process(peer *Peer) {
	// TODO: послать SteerMsg
	logrus.Println("state normal")
	if peer.step == 0 {
		// каждый Normal это VIEWER который умеет отправлять SteerMsg
		GetStateByRole(protocol.NodeRole_VIEWER, false).Process(peer)
	}
}
