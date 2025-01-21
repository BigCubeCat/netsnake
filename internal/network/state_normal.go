package network

import protocol "github.com/bigcubecat/netsnake/proto"

func (state PeerNormalState) Process(peer *Peer) {
	if peer.step == 0 {
		// каждый Normal это VIEWER который умеет отправлять SteerMsg
		GetStateByRole(protocol.NodeRole_NORMAL, false).Process(peer)
	}
}
