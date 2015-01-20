package main

import "fmt"

//Map Manager
func mapman(ch chan command) {
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
				m[cmd.key] = value{cmd.data, cmd.numbytes, version + 1, cmd.expiry}
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
					r = fmt.Sprintf("VALUE %v %v %v\r\n"+val.data+"\r\n", val.version, val.expiry, val.numbytes)
				}
			}
		case 3:
			{
				if ok {
					if val.version == cmd.version {
						m[cmd.key] = value{cmd.data, cmd.numbytes, val.version + 1, cmd.expiry}
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
		}
		cmd.resp <- r
	}
}
