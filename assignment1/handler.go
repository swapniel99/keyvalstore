package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func handleConn(conn net.Conn, ch chan<- *command) {
	addr := conn.RemoteAddr()
	log.Println(addr, "connected.")
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	resp := make(chan string, 1) // Response channel

	for {
		//Command Prompt
		//		write(writer, addr, "kv@cs733 ~ $ ")	// The Command Prompt :)

		str, e := reader.ReadString('\n')
		if e != nil {
			//Read error
			log.Println("ERROR reading:", addr, e)
			break
		}

		//Scan next line
		str = strings.TrimRight(str, "\r\n")
		if str == "" {
			continue //Empty command
		}

		cmd, e := parser(str)
		if e != nil {
			write(writer, addr, e.Error())
		} else {
			//Do work here
			cmd.resp = resp
			if cmd.action == 0 || cmd.action == 3 {
				buf := make([]byte, cmd.numbytes)
				_, ed := io.ReadFull(reader, buf)
				if (ed) != nil {
					//Read error
					log.Println("ERROR reading data:", addr, ed)
					break
				}
				tail, ed2 := reader.ReadString('\n')
				if (ed2) != nil {
					//Read error
					log.Println("ERROR reading post-data:", addr, ed2)
					break
				}
				cmd.data = string(buf)
				if (strings.TrimRight(tail,"\r\n") != "") || (len(cmd.data) != cmd.numbytes) {
					write(writer, addr, "ERR_CMD_ERR\r\n")
					continue
				}
			}
			ch <- &cmd
			reply := <-resp
			if cmd.action == 0 || cmd.action == 3 {
				if cmd.noreply {
					continue
				}
			}
			write(writer, addr, reply)
		}
	}
	// Shut down the connection.
	log.Println("Closing connection", addr)
	conn.Close()
}

//Writes in TCP connection
func write(w *bufio.Writer, a net.Addr, s string) {
	_, err := fmt.Fprintf(w, s)
	if err != nil {
		log.Println("ERROR writing:", a, err)
	}
	err = w.Flush()
	if err != nil {
		log.Println("ERROR flushing:", a, err)
	}
}
