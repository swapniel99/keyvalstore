package main

import (
	"log"
	"net"
)

type value struct {
	data     string
	numbytes int
	version  uint64
	expiry   int64
}

type command struct {
	action   int // 0=set, 1=get, 2=getm, 3=cas, 4=delete, [5=cleanup]-hidden
	key      string
	expiry   int64
	version  uint64
	numbytes int
	noreply  bool
	data     string
	resp     chan string
}

func main() {
	// Listen on TCP port 9000 on all interfaces.
	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	ch := make(chan command)
	go mapman(ch)
	go cleaner(1, ch)
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go handleConn(conn, ch)
	}
}
