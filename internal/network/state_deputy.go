package network

import (
	protocol "github.com/bigcubecat/netsnake/proto"
)

func (state PeerDeputyState) Process(peer *Peer) {
	// TODO: сделать проверку, что MASTER жив
	GetStateByRole(protocol.NodeRole_NORMAL).Process(peer)
}
