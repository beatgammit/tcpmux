// +build go1.3

package tcpmux

import (
	"net"
	"testing"
)

func BenchmarkNoMuxParallel(b *testing.B) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		b.Fatal(err)
	}

	defer l.Close()
	go serve(l)

	addr := l.Addr()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := connect(addr); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkMuxParallel(b *testing.B) {
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
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := connect(addr); err != nil {
				b.Fatal(err)
			}
		}
	})
}
