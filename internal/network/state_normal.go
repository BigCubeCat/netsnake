package network

import protocol "github.com/bigcubecat/netsnake/proto"

func (state PeerNormalState) Process(peer *Peer) {

	// каждый Normal это VIEWER который умеет отправлять SteerMsg
	GetStateByRole(protocol.NodeRole_NORMAL).Process(peer)
}
