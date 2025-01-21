package master

import (
	"errors"
	"net"
	"sync"

	"github.com/bigcubecat/netsnake/internal/network/peer"
	snakes "github.com/bigcubecat/netsnake/proto"
)

type Master struct {
	peers      map[int32]*peer.Peer
	nextPeerID int32
	mu         sync.Mutex
	gameConfig *snakes.GameConfig
}

func NewMaster(config *snakes.GameConfig) *Master {
	return &Master{
		peers:      make(map[int32]*peer.Peer),
		nextPeerID: 1,
		gameConfig: config,
	}
}

func (m *Master) HandleJoin(joinMsg *snakes.GameMessage_JoinMsg, addr *net.UDPAddr) (int32, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Проверка возможности присоединения
	if *joinMsg.RequestedRole == snakes.NodeRole_VIEWER {
		return m.addViewer(addr), nil
	}

	if !m.canAddPlayer() {
		return 0, errors.New("no space for new players")
	}

	// Создание нового игрока
	playerID := m.nextPeerID
	m.nextPeerID++

	m.peers[playerID] = &peer.Peer{
		ID:   playerID,
		Addr: addr,
		Role: *joinMsg.RequestedRole,
	}

	return playerID, nil
}

func (m *Master) canAddPlayer() bool {
	// Логика проверки доступного места на поле
	return true // Заглушка
}

func (m *Master) addViewer(addr *net.UDPAddr) int32 {
	viewerID := m.nextPeerID
	m.nextPeerID++
	m.peers[viewerID] = &peer.Peer{
		ID:   viewerID,
		Addr: addr,
		Role: snakes.NodeRole_VIEWER,
	}
	return viewerID
}
