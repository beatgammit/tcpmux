package tcpmux

import (
	"io/ioutil"
	"net"
	"testing"
)

func serve(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		ioutil.ReadAll(conn)
		conn.Close()
	}
}

func connect(addr net.Addr) error {
	conn, err := net.Dial(addr.Network(), addr.String())
	if err != nil {
		return err
	}
	conn.Write([]byte("Hello World"))
	return conn.Close()
}

func BenchmarkNoMuxSequential(b *testing.B) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		b.Fatal(err)
	}
	defer l.Close()

	go serve(l)

	addr := l.Addr()
	for i := 0; i < b.N; i++ {
		if err := connect(addr); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMuxSequential(b *testing.B) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		b.Fatal(err)
	}
	defer l.Close()

	go serve(l)

	lMux, err := net.Listen("tcp", ":0")
	if err != nil {
		b.Fatal(err)
	}
	defer lMux.Close()

	m := New(lMux, NetPipe{Network: l.Addr().Network(), Address: l.Addr().String()})
	go serve(m)

	addr := lMux.Addr()
	for i := 0; i < b.N; i++ {
		if err := connect(addr); err != nil {
			b.Fatal(err)
		}
	}
}
