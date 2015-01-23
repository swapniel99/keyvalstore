package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func handleConn(conn net.Conn, ch chan<- *command) {
	addr := conn.RemoteAddr()
	log.Println(addr, "connected.")
	scanner := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)
	resp := make(chan string, 1) // Response channel

	for {
		//Command Prompt
		//		write(writer, addr, "kv@cs733 ~ $ ")	// The Command Prompt :)

		success := scanner.Scan()
		if !success {
			e := scanner.Err()
			if e != nil {
				//Read error
				log.Println("ERROR reading:", addr, e)
			} else {
				//EOF received
				log.Println("End of Transmission by", addr)
			}
			break
		}

		//Scan next line
		str := scanner.Text()
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
				successd := scanner.Scan()
				if !successd {
					ed := scanner.Err()
					if ed != nil {
						//Read error
						log.Println("ERROR reading:", addr, ed)
					} else {
						//EOF received
						log.Println("End of Transmission by", addr)
					}
					break
				}
				cmd.data = scanner.Text()
				if len(cmd.data) != cmd.numbytes {
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
