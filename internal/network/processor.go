package network

import protocol "github.com/bigcubecat/netsnake/proto"

// интерфейс для обработки игры в зависимости от роли узла
type PeerProcessorState interface {
	Process(peer *Peer)
}

type PeerMasterState struct{}

type PeerDeputyState struct{}

type PeerNormalState struct{}

type PeerViewerState struct{}

type PeerJoinState struct{}

var states = [...]PeerProcessorState{
	PeerMasterState{},
	PeerDeputyState{},
	PeerNormalState{},
	PeerViewerState{},
	PeerJoinState{},
}

func GetStateByRole(role protocol.NodeRole, join bool) PeerProcessorState {
	if join {
		return states[4]
	}
	switch role {
	case protocol.NodeRole_MASTER:
		return states[0]
	case protocol.NodeRole_DEPUTY:
		return states[1]
	case protocol.NodeRole_NORMAL:
		return states[2]
	default:
		return states[3]
	}
}
