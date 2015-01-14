package tcpmux

import (
	"net"
)

type bufConn struct {
	*Conn
}

func (br *bufConn) Read(p []byte) (n int, err error) {
	n, err = br.Conn.Read(p)
	br.buf = append(br.buf, p...)
	return
}

type Conn struct {
	net.Conn

	buf []byte
}

func (r *Conn) Read(p []byte) (n int, err error) {
	if len(r.buf) > 0 {
		if len(p) < len(r.buf) {
			n = copy(p, r.buf)
			r.buf = r.buf[n:]
			return
		} else {
			n = copy(p, r.buf)
			r.buf = nil
		}
	}

	var n2 int
	n2, err = r.Conn.Read(p[n:])
	return n + n2, err
}
