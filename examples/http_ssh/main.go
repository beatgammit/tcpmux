package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/beatgammit/tcpmux"
)

var (
	addr     string
	httpAddr string
	sshAddr  string
)

func init() {
	flag.StringVar(&addr, "addr", ":8080", "address to listen to connections on")
	flag.StringVar(&sshAddr, "sshaddr", "localhost:22", "address to forward ssh connections to")
	flag.StringVar(&httpAddr, "httpaddr", "", "address to forward http connections to; leave blank to serve from current directory")
	flag.Parse()
}

func main() {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	protos := []tcpmux.Proto{
		tcpmux.SSH{tcpmux.NetPipe{Network: "tcp", Address: sshAddr}},
	}
	if httpAddr != "" {
		protos = append(protos, tcpmux.NetPipe{Network: "tcp", Address: httpAddr})
	}

	m := tcpmux.New(l, protos...)

	if httpAddr != "" {
		for {
			c, err := m.Accept()
			if err != nil {
				log.Fatalln(err)
			}
			log.Println("Unhandled connection:", c)
			c.Close()
		}
	} else {
		http.Handle("/", http.FileServer(http.Dir(".")))
		s := &http.Server{}
		s.Serve(m)
	}
}
