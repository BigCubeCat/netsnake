package network

import protocol "github.com/bigcubecat/netsnake/proto"

func (state PeerViewerState) Process(peer *Peer) {
	multicastInbox := peer.annoncementController.ReadInbox()
	for _, recvMessage := range multicastInbox {
		switch recvMessage.GetType().(type) {
		case *protocol.GameMessage_Announcement:
			// нас интересует только анансы
		}
	}
}
