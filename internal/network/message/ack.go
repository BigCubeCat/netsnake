package message

import (
	protocol "github.com/bigcubecat/netsnake/proto"
	"google.golang.org/protobuf/proto"
)

func NewAckMsg(msgSeq int64, senderId int32, receiverId int32) *protocol.GameMessage {
	return &protocol.GameMessage{
		MsgSeq:     proto.Int64(msgSeq),
		SenderId:   proto.Int32(senderId),
		ReceiverId: proto.Int32(receiverId),
		Type: &protocol.GameMessage_Ack{
			Ack: &protocol.GameMessage_AckMsg{},
		},
	}
}

func NewDiscoverMsg(msgSeq int64) *protocol.GameMessage {
	return &protocol.GameMessage{
		MsgSeq: proto.Int64(msgSeq),
		Type: &protocol.GameMessage_Discover{
			Discover: &protocol.GameMessage_DiscoverMsg{},
		},
	}
}

func NewPingMsg(msgSeq int64, senderId int32, receiverId int32) *protocol.GameMessage {
	return &protocol.GameMessage{
		MsgSeq:     proto.Int64(msgSeq),
		SenderId:   proto.Int32(senderId),
		ReceiverId: proto.Int32(receiverId),
		Type: &protocol.GameMessage_Ping{
			Ping: &protocol.GameMessage_PingMsg{},
		},
	}
}
