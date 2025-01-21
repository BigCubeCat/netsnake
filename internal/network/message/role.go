package message

import (
	protocol "github.com/bigcubecat/netsnake/proto"
	"google.golang.org/protobuf/proto"
)

func NewRoleChangeMsg(
	msgSeq int64,
	senderId int32,
	receiverId int32,
	senderRole *protocol.NodeRole,
	receiverRole *protocol.NodeRole,
) *protocol.GameMessage {
	return &protocol.GameMessage{
		MsgSeq:     proto.Int64(msgSeq),
		SenderId:   proto.Int32(senderId),
		ReceiverId: proto.Int32(receiverId),
		Type: &protocol.GameMessage_RoleChange{
			RoleChange: &protocol.GameMessage_RoleChangeMsg{
				SenderRole:   senderRole,
				ReceiverRole: receiverRole,
			},
		},
	}
}
