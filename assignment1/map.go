package main

import (
	"fmt"
	//	"log"
	"time"
)

//Map Manager
func mapman(ch chan *command) {
	//The map which actually stores values
	m := make(map[string]value)
	go cleaner(1, ch)
	for cmd := range ch {
		r := "ERR_NOT_FOUND\r\n"
		val, ok := m[cmd.key]
		switch cmd.action {
		case 0:
			{
				var version uint64
				if !ok {
					version = 0
				} else {
					version = val.version
				}
				t := cmd.expiry
				if t != 0 {
					t += time.Now().Unix()
				}
				m[cmd.key] = value{cmd.data, cmd.numbytes, version + 1, t}
				r = fmt.Sprintf("OK %v\r\n", version+1)
			}
		case 1:
			{
				if ok {
					r = fmt.Sprintf("VALUE %v\r\n"+val.data+"\r\n", val.numbytes)
				}
			}
		case 2:
			{
				if ok {
					t := val.expiry
					if t != 0 {
						t = val.expiry - time.Now().Unix() // remaining time
					}
					if t < 0 {
						t = 0
					}
					r = fmt.Sprintf("VALUE %v %v %v\r\n"+val.data+"\r\n", val.version, t, val.numbytes)
				}
			}
		case 3:
			{
				if ok {
					if val.version == cmd.version {
						t := cmd.expiry
						if t != 0 {
							t += time.Now().Unix()
						}
						m[cmd.key] = value{cmd.data, cmd.numbytes, val.version + 1, t}
						r = fmt.Sprintf("OK %v\r\n", val.version+1)
					} else {
						r = fmt.Sprintf("ERR_VERSION\r\n")
					}
				}
			}
		case 4:
			{
				if ok {
					delete(m, cmd.key)
					r = "DELETED\r\n"
				}
			}
		case 5:
			{
				for k, v := range m {
					if v.expiry != 0 && v.expiry-time.Now().Unix() < 0 {
						delete(m, k)
					}
				}
				r = "CLEANED\r\n"
			}
		default:
			{
				r = "ERR_INTERNAL\r\n"
			}
		}
		cmd.resp <- r
	}
}

func cleaner(interval int, ch chan<- *command) {
	resp := make(chan string, 2) // Response channel
	c := command{5, "", 0, 0, 0, false, "", resp}
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		ch <- &c
		<-resp //	Receive "CLEANED\r\n" message
		//		log.Println(s)
	}
}
