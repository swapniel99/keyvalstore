package main

import (
	"errors"
	"strconv"
	"strings"
)

func parser(cmd string) (command, error) {
	arr := strings.Split(cmd, " ")
	c := command{0, "", 0, 0, 0, false, "", nil}
	e := errors.New("ERR_CMD_ERR\r\n")
	l := len(arr)
	switch arr[0] {
	case "set":
		{
			if l != 4 && l != 5 {
				return c, e
			}
			c.action = 0
			c.key = arr[1]
			exp, e1 := strconv.Atoi(arr[2])
			if e1 != nil || exp < 0 {
				return c, e
			}
			c.expiry = int64(exp)
			numb, e2 := strconv.Atoi(arr[3])
			if e2 != nil || numb < 0 {
				return c, e
			}
			c.numbytes = numb
			if l == 5 {
				if arr[4] == "noreply" {
					c.noreply = true
				} else {
					return c, e
				}
			}
			return c, nil
		}
	case "get":
		{
			if l != 2 {
				return c, e
			}
			c.action = 1
			c.key = arr[1]
			return c, nil
		}
	case "getm":
		{
			if l != 2 {
				return c, e
			}
			c.action = 2
			c.key = arr[1]
			return c, nil
		}
	case "cas":
		{
			if l != 5 && l != 6 {
				return c, e
			}
			c.action = 3
			c.key = arr[1]
			exp, e1 := strconv.Atoi(arr[2])
			if e1 != nil || exp < 0 {
				return c, e
			}
			c.expiry = int64(exp)
			ver, e2 := strconv.Atoi(arr[3])
			if e2 != nil || ver <= 0 {
				return c, e
			}
			c.version = uint64(ver)
			numb, e3 := strconv.Atoi(arr[4])
			if e3 != nil || numb < 0 {
				return c, e
			}
			c.numbytes = numb
			if l == 6 {
				if arr[5] == "noreply" {
					c.noreply = true
				} else {
					return c, e
				}
			}
			return c, nil
		}
	case "delete":
		{
			if l != 2 {
				return c, e
			}
			c.action = 4
			c.key = arr[1]
			return c, nil
		}
	case "cleanup": // Not specified in syntax, but provides manual cleanup option
		{
			if l != 1 {
				return c, e
			}
			c.action = 5
			return c, nil
		}
	default:
		{
			return c, e
		}
	}
}
