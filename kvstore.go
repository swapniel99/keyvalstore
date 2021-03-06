package main

import (
	"log"
	"net"
	"time"
)

type NullWriter int

func (NullWriter) Write([]byte) (int, error) {
	return 0, nil
}

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

	log.SetOutput(new(NullWriter))

	var l net.Listener
	var err error

	for {
		l, err = net.Listen("tcp", ":9000")
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second)
		} else {
			break
		}
	}

	defer l.Close()
	ch := make(chan *command, 10000)
	go mapman(ch)
	for {
		// Wait for a connection.
		conn, err2 := l.Accept()
		if err2 != nil {
			log.Println(err)
			continue
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go handleConn(conn, ch)
	}
}
