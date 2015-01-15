package tcpmux

import (
	"net"
)

type NoHandler struct {
	Conn *Conn
}

func (err NoHandler) Error() string {
	return "Unhandled connection"
}

type Proto interface {
	// Matches returns whether this net.Conn matches this protocol.
	//
	// Calls to Read() are buffered, so no data is consumed from the underlying socket.
	Matches(net.Conn) bool
	// Handle takes control of this net.Conn. The net.Conn may not have matched this protocol
	// if this protocol was set as the default.
	Handle(net.Conn) error
}

type Mux struct {
	Default Proto

	l      net.Listener
	protos []Proto
}

func New(l net.Listener, protos ...Proto) *Mux {
	return &Mux{l: l, protos: protos}
}

func (m *Mux) Accept() (net.Conn, error) {
	for {
		c, err := m.l.Accept()
		if err != nil {
			return nil, err
		}
		if err = m.Handle(c); err != nil {
			if c, ok := err.(NoHandler); ok {
				return c.Conn, nil
			}
			// TODO: handle other errors
		}
	}
}

func (m *Mux) Close() error {
	return m.l.Close()
}

func (m *Mux) Addr() net.Addr {
	return m.l.Addr()
}

func (m *Mux) Handle(c net.Conn) error {
	buf := &bufConn{Conn: &Conn{Conn: c}}
	for _, p := range m.protos {
		if p.Matches(buf) {
			return p.Handle(buf.Conn)
		}
		buf.Reset()
	}
	if m.Default != nil {
		return m.Default.Handle(buf.Conn)
	}
	return NoHandler{buf.Conn}
}
