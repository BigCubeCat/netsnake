package announcementcontroller

import (
	"fmt"
	"net"

	"github.com/bigcubecat/netsnake/internal/config"
)

// InboxRoutine запускает прослушивание multicast адреса
func InboxRoutine(config config.CliConfig) {
	addr, err := net.ResolveUDPAddr("udp", config.MulticastAddress)
	if err != nil {
		panic("error resolve multicast address " + config.MulticastAddress)
	}
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("error listening on multicast:", err)
		return
	}
	defer conn.Close()

}
