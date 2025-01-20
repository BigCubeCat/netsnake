package message

import (
	"errors"
	snakes "path/to/your/generated/proto"
	"time"

	"github.com/golang/protobuf/proto"
)

func EncodeGameMessage(msg *snakes.GameMessage) ([]byte, error) {
	return proto.Marshal(msg)
}

func DecodeGameMessage(data []byte) (*snakes.GameMessage, error) {
	msg := &snakes.GameMessage{}
	err := proto.Unmarshal(data, msg)
	return msg, err
}

func CreateJoinMessage(name, gameName string, role snakes.NodeRole) *snakes.GameMessage {
	return &snakes.GameMessage{
		MsgSeq: proto.Int64(time.Now().UnixNano()),
		Type: &snakes.GameMessage_Join{
			Join: &snakes.GameMessage_JoinMsg{
				PlayerName:    proto.String(name),
				GameName:      proto.String(gameName),
				RequestedRole: &role,
			},
		},
	}
}
