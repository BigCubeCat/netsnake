package message

import (
	protocol "github.com/bigcubecat/netsnake/proto"

	"google.golang.org/protobuf/proto"
)

// MarshalGameMessage кодирует Сообщение в бинарный объект
func MarshalGameMessage(msg *protocol.GameMessage) ([]byte, error) {
	return proto.Marshal(msg)
}

// UnmarshalGameMessage декодирует Сообщение из бинарного объекта
func UnmarshalGameMessage(data []byte) (*protocol.GameMessage, error) {
	msg := &protocol.GameMessage{}
	err := proto.Unmarshal(data, msg)
	return msg, err
}
