package tcpmux

import (
	"bytes"
	"io"
	"math/rand"
	"net"
	"testing"
	"time"
)

type fakeConn struct {
	io.Reader
	closed bool
}

func (conn *fakeConn) Read(b []byte) (int, error) {
	if conn.closed {
		return 0, io.EOF
	}
	return conn.Reader.Read(b)
}
func (conn *fakeConn) Write(b []byte) (int, error)        { return len(b), nil }
func (conn *fakeConn) Close() error                       { conn.closed = true; return nil }
func (conn *fakeConn) LocalAddr() net.Addr                { return nil }
func (conn *fakeConn) RemoteAddr() net.Addr               { return nil }
func (conn *fakeConn) SetDeadline(t time.Time) error      { return nil }
func (conn *fakeConn) SetReadDeadline(t time.Time) error  { return nil }
func (conn *fakeConn) SetWriteDeadline(t time.Time) error { return nil }

func TestConn(t *testing.T) {
	const (
		start = "Hello"
		read  = "World"
	)

	conn := &Conn{
		Conn: &fakeConn{Reader: bytes.NewBufferString(read)},
		buf:  []byte(start),
	}
	buf := make([]byte, 5)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("Read %d bytes, but error: %s", n, err)
	}

	if string(buf) != start {
		t.Errorf("Didn't read from buf correctly: %s != %s", string(buf), start)
	}
	if len(conn.buf) != 0 {
		t.Error("buf not cleared after reading")
	}

	buf = make([]byte, 5)
	_, err = conn.Read(buf)
	if err != nil {
		t.Fatal(err)
	}

	if string(buf) != read {
		t.Errorf("Didn't read from Reader correctly: %s != %s", string(buf), read)
	}

	_, err = conn.Read(buf)
	if err == nil {
		t.Error("Expected error when socket should be empty")
	}
}

func TestBufConnReset(t *testing.T) {
	const (
		str = "Hello world"
	)

	conn := &bufConn{
		Conn: &Conn{
			Conn: &fakeConn{Reader: bytes.NewBufferString(str)},
		},
	}

	buf := make([]byte, len(str))
	// 73 is the best number and is used to keep results consistent
	rand.Seed(73)
	for i := 0; i < 100; i++ {
		n := rand.Intn(len(str))
		n2, err := conn.Read(buf[:n])
		if err != nil {
			t.Fatal(err)
		}

		if n2 != n {
			t.Errorf("%d != %d", n2, n)
		}
		if string(buf[:n]) != str[:n] {
			t.Errorf("%d: %s != %s", n, string(buf[:n]), str[:n])
		}
		conn.Reset()
	}
}
