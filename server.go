package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type key string

type value struct {
	data     string
	numbytes uint64
	version  uint64
	expiry   uint64
}

type command struct {
	action   int // 0=set, 1=get, 2=getm, 3=cas, 4=delete
	id       key
	expiry   uint64
	version  uint64
	numbytes uint64
	noreply  bool
	data     string
	resp     chan response
}

type response struct {
	msg  string
	data string
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

func handleConn(conn net.Conn, ch chan command) {
	// Echo all incoming data.
	log.Println(conn.RemoteAddr(), "connected.")
	scanner := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)
	resp := make(chan response)
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
		//		str := "world"
		ch <- command{0, "hello", 0, 0, uint64(len(str)), true, str, resp}
		r := <-resp
		fmt.Println(r.msg, r.data)
		ch <- command{1, "hello", 0, 0, 0, true, "", resp}
		r = <-resp
		fmt.Print(r.msg, r.data)

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

//Map Manager
func mapman(ch chan command) {
	//The map which actually stores values
	m := make(map[key]value)
	for cmd := range ch {
		val, ok := m[cmd.id]
		switch cmd.action {
		case 0:
			{
				var version uint64
				if !ok {
					version = 0
				} else {
					version = val.version
				}
				m[cmd.id] = value{cmd.data, cmd.numbytes, version + 1, cmd.expiry}
				cmd.resp <- response{fmt.Sprintf("OK %v\r\n", version+1), ""}
			}
		case 1:
			{
				if !ok {
					cmd.resp <- response{fmt.Sprintf("ERR_NOT_FOUND\r\n"), ""}
				} else {
					cmd.resp <- response{fmt.Sprintf("VALUE %v\r\n", val.numbytes), val.data}
				}
			}
		case 2:
			{
				if !ok {
					cmd.resp <- response{fmt.Sprintf("ERR_NOT_FOUND\r\n"), ""}
				} else {
					cmd.resp <- response{fmt.Sprintf("VALUE %v %v %v\r\n", val.version, val.expiry, val.numbytes), val.data}
				}
			}
		case 3:
			{
				if !ok {
					cmd.resp <- response{fmt.Sprintf("ERR_NOT_FOUND\r\n"), ""}
				} else {
					if val.version == cmd.version {
						m[cmd.id] = value{cmd.data, cmd.numbytes, val.version + 1, cmd.expiry}
						cmd.resp <- response{fmt.Sprintf("OK %v\r\n", val.version+1), ""}
					} else {
						cmd.resp <- response{fmt.Sprintf("ERR_VERSION\r\n"), ""}
					}
				}
			}
		case 4:
			{
				if !ok {
					cmd.resp <- response{fmt.Sprintf("ERR_NOT_FOUND\r\n"), ""}
				} else {
					delete(m, cmd.id)
					cmd.resp <- response{fmt.Sprintf("DELETED\r\n"), ""}
				}
			}
		}
	}
}
