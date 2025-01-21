package message

import (
	"time"

	protocol "github.com/bigcubecat/netsnake/proto"
	"google.golang.org/protobuf/proto"
)

func NewAnnouncementMessage(
	name, gameName string,
	width, height, foodStatic, stateDelay int32,
	role protocol.NodeRole,
	players *protocol.GamePlayers,
) *protocol.GameMessage {
	return &protocol.GameMessage{
		MsgSeq: proto.Int64(time.Now().UnixNano()),
		Type: &protocol.GameMessage_Announcement{
			Announcement: &protocol.GameMessage_AnnouncementMsg{
				Games: []*protocol.GameAnnouncement{
					{
						GameName: proto.String(gameName),
						Config: &protocol.GameConfig{
							Width:        proto.Int32(width),
							Height:       proto.Int32(height),
							FoodStatic:   proto.Int32(foodStatic),
							StateDelayMs: proto.Int32(stateDelay),
						},
						CanJoin: proto.Bool(true),
						Players: players,
					},
				},
			},
		},
	}
}
