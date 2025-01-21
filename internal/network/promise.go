package network

import protocol "github.com/bigcubecat/netsnake/proto"

type MessagePromise struct {
	Address string
	Port    int
	Message *protocol.GameMessage
}
