package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	tcs "github.com/yvv4git/tcpserver"
)

var clients []*tcs.Client

func main() {
	server := tcs.NewServer("127.0.0.1:1234")

	// add new client to slice
	server.OnNewClient(func(c *tcs.Client) {
		fmt.Printf("Attach client: [%p] \n", c)
		clients = append(clients, c)
		printEnter()
	})

	// print message from client
	server.OnNewMessage(func(c *tcs.Client, message string) {
		fmt.Printf("Message[%p]:\n%s", c, message)
		printEnter()
	})

	// delete client from slice if his detach
	server.OnClientConnectionClosed(func(c *tcs.Client, err error) {
		fmt.Printf("Detach client: [%p] \n", c)
		for key, val := range clients {
			if val == c {
				clients = append(clients[:key], clients[key+1:]...)
				break
			}
		}
	})

	// stand up server in new gorutine
	go server.Listen()

	// read admin commands
	cmdReader := bufio.NewReader(os.Stdout)
	for {
		printEnter()
		cmd, _ := cmdReader.ReadString('\n')
		cmdProcessing(cmd)
	}
}

// processing admin command
func cmdProcessing(cmd string) {
	fmt.Printf("Cmd: %s\n", cmd)
	reCmd := regexp.MustCompile(`^(\d+):`) // приставка, которая означает, что дальше идет комманда для ОС

	if cmd == "ls\n" {
		fmt.Println("===clients===")
		for key, val := range clients {
			fmt.Printf("[%d][%p]\n", key, val)
		}
	} else if cmd == "exit\n" {
		fmt.Println("By-by")
		os.Exit(1)
	} else if reCmd.MatchString(cmd) {
		cmdExec := reCmd.ReplaceAllString(cmd, "") // парсим комманду

		// парсим порядковый номер клиента
		botId := reCmd.FindString(cmd)
		botId = strings.TrimSuffix(botId, ":")
		botID, _ := strconv.Atoi(botId)

		if botID < len(clients) {
			// шлем сообщение клиенту
			clients[botID].Send(cmdExec)
		}
	}
}

func printEnter() {
	fmt.Println("\nEnter command: ")
}
