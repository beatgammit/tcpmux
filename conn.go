package tcpmux

import (
	"net"
)

type bufConn struct {
	*Conn

	buf []byte
}

func (br *bufConn) Reset() {
	br.Conn.buf = append(br.buf, br.Conn.buf...)
	br.buf = nil
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
		n = copy(p, r.buf)
		if n <= len(r.buf) {
			r.buf = r.buf[n:]
		} else {
			r.buf = nil
		}
	}

	var n2 int
	n2, err = r.Conn.Read(p[n:])
	return n + n2, err
}
