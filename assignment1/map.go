package main

import (
	"fmt"
	"time"
)

//Map Manager
func mapman(ch <-chan command) {
	//The map which actually stores values
	m := make(map[string]value)
	for cmd := range ch {
		val, ok := m[cmd.key]
		r := "ERR_NOT_FOUND\r\n"
		switch cmd.action {
		case 0:
			{
				var version uint64
				if !ok {
					version = 0
				} else {
					version = val.version
				}
				m[cmd.key] = value{cmd.data, cmd.numbytes, version + 1, time.Now().Unix() + cmd.expiry}
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
					t := val.expiry - time.Now().Unix() // remaining time
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
						m[cmd.key] = value{cmd.data, cmd.numbytes, val.version + 1, time.Now().Unix() + cmd.expiry}
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
					if v.expiry-time.Now().Unix() < 0 {
						delete(m, k)
					}
				}
				r = "CLEANED\r\n"
			}
		}
		cmd.resp <- r
	}
}

func cleaner(interval int, ch chan<- command) {
	resp := make(chan string) // Response channel
	c := command{5, "", 0, 0, 0, false, "", resp}
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		ch <- c
		s := <-resp
		fmt.Println(s)
	}
}
