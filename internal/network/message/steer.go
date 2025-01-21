package message

import (
	protocol "github.com/bigcubecat/netsnake/proto"
	"google.golang.org/protobuf/proto"
)

// Создать StateMsg - сообщение с расположением игровых объектов
func NewStateMsg(
	msgSeq int64,
	senderId int32,
	receiverId int32,
	stateOrder int32,
	snakes []*protocol.GameState_Snake,
	foods []*protocol.GameState_Coord,
	players *protocol.GamePlayers,
) *protocol.GameMessage {
	return &protocol.GameMessage{
		MsgSeq:     proto.Int64(msgSeq),
		SenderId:   proto.Int32(senderId),
		ReceiverId: proto.Int32(receiverId),
		Type: &protocol.GameMessage_State{
			State: &protocol.GameMessage_StateMsg{
				State: &protocol.GameState{
					StateOrder: proto.Int32(stateOrder),
					Snakes:     snakes,
					Foods:      foods,
					Players:    players,
				},
			},
		},
	}
}

// Создать SteerMsg - сообщение для поворота своей змейки
func NewSteerMsg(
	msgSeq int64,
	senderId int32,
	receiverId int32,
	direction int32,
) *protocol.GameMessage {
	return &protocol.GameMessage{
		MsgSeq:     proto.Int64(msgSeq),
		SenderId:   proto.Int32(senderId),
		ReceiverId: proto.Int32(receiverId),
		Type: &protocol.GameMessage_Steer{
			Steer: &protocol.GameMessage_SteerMsg{
				Direction: (*protocol.Direction)(proto.Int32((int32)(direction))),
			},
		},
	}
}
