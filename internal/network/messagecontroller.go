package network

import (
	"context"
	"net"

	protocol "github.com/bigcubecat/netsnake/proto"
)

type MessageController struct {
	conn    *net.UDPConn
	ctx     *context.Context
	peerPtr *Peer

	InboxMessages []*protocol.GameMessage
}

func NewMessageController(ctx *context.Context, peerPtr *Peer) *MessageController {
	return &MessageController{
		ctx:     ctx,
		peerPtr: peerPtr,
	}
}
