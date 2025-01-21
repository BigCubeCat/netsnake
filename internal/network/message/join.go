package message

import (
	protocol "github.com/bigcubecat/netsnake/proto"

	"time"

	"google.golang.org/protobuf/proto"
)

func CreateJoinMessage(name, gameName string, role protocol.NodeRole) *protocol.GameMessage {
	return &protocol.GameMessage{
		MsgSeq: proto.Int64(time.Now().UnixNano()),
		Type: &protocol.GameMessage_Join{
			Join: &protocol.GameMessage_JoinMsg{
				PlayerName:    proto.String(name),
				GameName:      proto.String(gameName),
				RequestedRole: &role,
			},
		},
	}
}
