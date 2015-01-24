package main

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestConcurrentSets(t *testing.T) {
	N := 10000 //Number of concurrent writes: 10 thousand
	ch := make(chan *command)
	resp := make(chan string, 1) // Response channel
	c := command{0, "key", 0, 0, 5, false, "value", resp}

	go mapman(ch)

	ack := make(chan bool, N)

	for i := 0; i < N; i++ {
		go func() {
			ch <- &c
			r := strings.Split(<-resp, " ")[0]
			if "OK" == r {
				ack <- true
			}
		}()
	}

	tick := time.Tick(time.Second)

	for i := 0; i < N; i++ {
		select {
		case <-tick:
			{
				t.Error("Timeout", N-i, "thread(s) did not return ack.")
				break
			}
		case <-ack:
		}
	}

	c.action = 2
	c.data = ""

	ch <- &c
	r := strings.Split(<-resp, " ")[1]
	if strconv.Itoa(N-1) != r {
		t.Error("Serial version mismatch : Got "+r+" instead of", N)
	}
}

func TestConcurrentGets(t *testing.T) {
	N := 10000 //Number of concurrent reads: 10 thousand
	ch := make(chan *command)
	resp := make(chan string, 1) // Response channel
	c := command{0, "key", 0, 0, 5, false, "value", resp}

	go mapman(ch)
	ch <- &c

	if "OK" != strings.Split(<-resp, " ")[0] {
		t.Error("Unable to set : ")
	}

	c.action = 1
	c.data = ""

	ack := make(chan bool, N)

	for i := 0; i < N; i++ {
		go func() {
			ch <- &c
			if "VALUE 5\r\nvalue\r\n" == <-resp {
				ack <- true
			}
		}()
	}

	tick := time.Tick(time.Second)

	for i := 0; i < N; i++ {
		select {
		case <-tick:
			{
				t.Error("Timeout", N-i, "thread(s) did not return ack.")
				break
			}
		case <-ack:
		}
	}
}

func TestExpiry(t *testing.T) {
	ch := make(chan *command)
	resp := make(chan string, 1) // Response channel
	c := command{0, "key", 2, 0, 5, false, "value", resp}

	go mapman(ch)
	ch <- &c
	if "OK" != strings.Split(<-resp, " ")[0] {
		t.Error("Unable to set.")
	}

	time.Sleep(5 * time.Second)

	c.action = 1
	c.data = ""
	c.expiry = 0

	ch <- &c
	r := <-resp
	if "ERR_NOT_FOUND\r\n" != r {
		t.Error("Value not expired : " + r)
	}
}

func TestConcurrentTCPSets(t *testing.T) {
	go main()

	N := 200 //Number of concurrent connections: 2 hundred
	k := 100 //Number of operations per thread
	//Total 200x100 = 20000 writes
	ack := make(chan bool, N)

	for i := 0; i < N; i++ {
		go client(t, k, ack, i)
	}

	tick := time.Tick(3 * time.Second)

	for i := 0; i < N; i++ {
		select {
		case <-tick:
			{
				t.Error("Timeout", N-i, "thread(s) did not return ack.")
				break
			}
		case <-ack:
		}
	}
}

func client(t *testing.T, k int, ack chan<- bool, id int) {
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		t.Error(err, id)
		ack <- true
		return
	}
	buff := bufio.NewReader(conn)
	var line []byte
	for j := 1; j < k; j++ {
		io.Copy(conn, bytes.NewBufferString("set key 0 9\r\nSomeValue\r\n"))

		line, err = buff.ReadBytes('\n')
		if err != nil {
			t.Error("Some error in reading TCP :", err)
			conn.Close()
			ack <- true
			return
		}
		input := strings.TrimRight(string(line), "\r\n")
		result := strings.Split(input, " ")[0]
		if result != "OK" {
			t.Error("Unable to set")
		}
	}
	conn.Close()
	ack <- true
}
