package message

import (
	"time"

	protocol "github.com/bigcubecat/netsnake/proto"
	"google.golang.org/protobuf/proto"
)

func CreateAnnouncementMessage(
	name, gameName string,
	role protocol.NodeRole,
) *protocol.GameMessage {
	return &protocol.GameMessage{
		MsgSeq: proto.Int64(time.Now().UnixNano()),
		Type: &protocol.GameMessage_Announcement{
			Announcement: &protocol.GameMessage_AnnouncementMsg{
			  Games: []*protocol.GameAnnouncement{
			    {
			    }
			  },
			},
		},
	}
}
