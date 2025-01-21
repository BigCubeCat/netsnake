package network

import (
	"github.com/bigcubecat/netsnake/internal/common"
	"github.com/bigcubecat/netsnake/internal/model"
	"github.com/bigcubecat/netsnake/internal/network/message"
	"github.com/bigcubecat/netsnake/internal/utils"
	protocol "github.com/bigcubecat/netsnake/proto"
)

func (peer *Peer) MoveSnake(direction int) {
	if peer.Role != protocol.NodeRole_VIEWER {
		peer.GameInstance.MoveSnake(peer.ID, direction)
	}
}

func (peer *Peer) JoinPlayer(joinMsg *protocol.GameMessage_JoinMsg, address string, port int) {
	id := utils.RandomId()
	playerName := joinMsg.GetPlayerName()
	role := joinMsg.GetRequestedRole()
	for _, v := range peer.Players {
		if v.Name == playerName {
			return
		}
	}
	peer.Players[id] = common.Player{
		ID:        id,
		Role:      roleToMode(role),
		Name:      playerName,
		IpAddress: address,
		Port:      port,
		Score:     0,
	}
	if role != protocol.NodeRole_VIEWER {
		peer.GameInstance.AddSnake(id, model.Normal)
	}
	peer.messageController.AddMessage(
		address,
		port,
		message.NewAckMsg(peer.msgSeq.Load(), int32(peer.ID), int32(id)),
	)
}
