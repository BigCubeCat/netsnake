package network_test

import (
	"net"
	"testing"

	"github.com/bigcubecat/netsnake/internal/network"
	network_message "github.com/bigcubecat/netsnake/internal/network/message"
	snakes "github.com/bigcubecat/netsnake/proto"
	"google.golang.org/protobuf/proto"
)

func TestJoinHandler(t *testing.T) {
	config := &snakes.GameConfig{
		Width:        proto.Int32(40),
		Height:       proto.Int32(30),
		FoodStatic:   proto.Int32(1),
		StateDelayMs: proto.Int32(1000),
	}

	n := network.NewNetwork(config)
	defer n.Conn.Close()

	// Тест успешного присоединения
	joinMsg := network_message.CreateJoinMessage("test", "game", snakes.NodeRole_NORMAL)
	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	n.HandleJoinRequest(joinMsg.GetJoin(), addr)

	if len(n.Master.Peers) != 1 {
		t.Errorf("Expected 1 peer, got %d", len(n.Master.peers))
	}
}
