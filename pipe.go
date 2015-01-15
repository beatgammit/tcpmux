package tcpmux

import (
	"io"
	"net"
)

type NetPipe struct {
	Network, Address string
}

func pipe(dst io.WriteCloser, src io.ReadCloser) error {
	_, err := io.Copy(dst, src)
	dst.Close()
	src.Close()
	return err
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
	go pipe(conn, c)
	go pipe(c, conn)
	return nil
}
