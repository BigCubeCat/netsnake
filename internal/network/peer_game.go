package network

import (
	protocol "github.com/bigcubecat/netsnake/proto"
)

func (peer *Peer) MoveSnake(direction int) {
	if peer.Role != protocol.NodeRole_VIEWER {
		peer.GameInstance.MoveSnake(peer.ID, direction)
	}
}
