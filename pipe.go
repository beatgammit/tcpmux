package tcpmux

import (
	"io"
	"net"
)

type NetPipe struct {
	Network, Address string
}

func (p NetPipe) Matches(c net.Conn) bool {
	return true
}
func (p NetPipe) Handle(c net.Conn) error {
	conn, err := net.Dial(p.Network, p.Address)
	if err != nil {
		return err
	}

	// TODO: handle errors
	go io.Copy(conn, c)
	go io.Copy(c, conn)
	return nil
}
