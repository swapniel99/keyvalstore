package main

import (
	"container/heap"
	"fmt"
	"time"
)

//Map Manager
func mapman(ch chan *command) {
	//The map which actually stores values
	m := make(map[string]value)
	h := &nodeHeap{}
	var counter uint64 = 0
	go cleaner(1, ch)
	for cmd := range ch {
		r := "ERR_NOT_FOUND\r\n"
		val, ok := m[cmd.key]
		switch cmd.action {
		case 0:
			{
				version := counter
				counter++
				t := cmd.expiry
				if t != 0 {
					t += time.Now().Unix()
				}
				m[cmd.key] = value{cmd.data, cmd.numbytes, version, t}
				if cmd.expiry != 0 {
					heap.Push(h, node{t, cmd.key, version})
				}
				r = fmt.Sprintf("OK %v\r\n", version)
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
						version := counter
						counter++
						m[cmd.key] = value{cmd.data, cmd.numbytes, version, t}
						if cmd.expiry != 0 {
							heap.Push(h, node{t, cmd.key, version})
						}
						r = fmt.Sprintf("OK %v\r\n", version)
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
				t := time.Now().Unix()
				for (*h).Len() != 0 && (*h)[0].expiry <= t {
					root := heap.Pop(h).(node)
					v, e := m[root.key]
					if e && root.version == v.version {
						delete(m, root.key)
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
	resp := make(chan string, 1) // Response channel
	c := command{5, "", 0, 0, 0, false, "", resp}
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		ch <- &c
		<-resp
	}
}
