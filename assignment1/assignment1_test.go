package main

import (
    "testing"
    "time"
)

func TestConcurrentRead(t *testing.T) {
	n:=10000	//Number of concurrent threads
	ch := make(chan command)
	resp := make(chan string) // Response channel
	c := command{0, "key", 0, 0, 5, false, "value", resp}
	
	go mapman(ch)
	ch<-c
	_=<-resp
	
	c.action=1
	c.data=""
	
	ack:=make(chan bool,n)
	
	for i:=0;i<n;i++ {
		go func() {
			ch<-c
			r:=<-resp
			if r=="VALUE 5\r\nvalue\r\n" {
				ack<-true
			}
		}()
	}
	
	tick := time.Tick(1*time.Second)
	
	for i:=0;i<n;i++ {
		select {
			case <-ack:
			case <-tick:
			{
				t.Error("Timeout")
				break
			}
		}
	}
}