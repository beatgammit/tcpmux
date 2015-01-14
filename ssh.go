package tcpmux

import (
	"net"
)

type SSH struct {
	NetPipe
}

func (p SSH) Matches(c net.Conn) bool {
	arr := make([]byte, 3)
	n, err := c.Read(arr)
	if err != nil || n != len(arr) {
		return false
	}
	return string(arr) == "SSH"
}

func (p SSH) Handle(c net.Conn) error {
	if p.Network == "" {
		p.Network = "tcp"
	}
	if p.Address == "" {
		p.Address = "localhost:22"
	}
	return p.NetPipe.Handle(c)
}
