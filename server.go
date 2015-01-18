package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	// Listen on TCP port 9000 on all interfaces.
	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	ch := make(chan string)
	go dowork(ch)
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

func handleConn(conn net.Conn, ch chan string) {
	// Echo all incoming data.
	log.Println(conn.RemoteAddr(), "connected.")
	scanner := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)
	for {
		success := scanner.Scan()
		if !success {
			e := scanner.Err()
			if e != nil {
				//Read error
				log.Println("ERROR reading:", conn.RemoteAddr(), e)
			} else {
				//EOF received
				log.Println("End of Transmission by", conn.RemoteAddr())
			}
			break
		}
		//Scan next line
		str := scanner.Text()

		//Do work here
		ch <- str

		//Send response
		_, err := fmt.Fprintf(writer, str+"\r\n")
		if err != nil {
			log.Println("ERROR writing:", conn.RemoteAddr(), err)
		}
		err = writer.Flush()
		if err != nil {
			log.Println("ERROR flushing:", conn.RemoteAddr(), err)
		}
	}
	// Shut down the connection.
	log.Println("Closing connection", conn.RemoteAddr())
	conn.Close()
}

func dowork(ch chan string) {
	//The map which actually stores values
	for s := range ch {
		fmt.Println(s)
	}
}
