package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type client struct {
	connection net.Conn
	username   string
}

var clients []*client

const host, port = "", "1337"

var allClients = 0

func main() {
	psock, err := net.Listen("tcp", host+":"+port)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Server started at %v:%v\n", host, port)

	for {
		conn, _ := psock.Accept()

		client := client{connection: conn, username: "User" + strconv.Itoa(allClients)}
		clients = append(clients, &client)
		fmt.Printf("Client (%v) connected\n", client.username)
		allClients++

		channel := make(chan string)
		go ProcessInput(channel, &client)
		go ProcessOutput(channel, &client)
	}
}

// ProcessInput handles all input from clients
func ProcessInput(out chan string, client *client) {
	defer close(out)

	reader := bufio.NewReader(client.connection)

	for {
		line, error := reader.ReadBytes('\n')
		message := strings.TrimSuffix(string(line), "\n")

		if error != nil {
			fmt.Printf("Client (%v) disconnected\n", client.username)
			client.connection.Close()
			clients = removeClient(clients, client)
			out <- "leaved"
			return
		}

		if message != "" {
			if message[0] == '/' {
				switch message[1:9] {
				case "username":
					client.username = message[10:]
					out <- "changed username"
				}
			} else {
				out <- message
			}
		}
	}
}

func removeClient(arr []*client, item *client) []*client {
	rtn := arr
	index := -1
	for i, value := range arr {
		if value == item {
			index = i
			break
		}
	}

	if index >= 0 {
		rtn = make([]*client, len(arr)-1)
		copy(rtn, arr[:index])
		copy(rtn[index:], arr[index+1:])
	}

	return rtn
}

// ProcessOutput outputs messages to all clients
func ProcessOutput(in <-chan string, client *client) {
	for {
		message := <-in
		if message != "" {
			for _, _client := range clients {
				if _client != client {
					fmt.Fprintln(_client.connection, client.username+": "+message)
				}
			}
		}
	}
}
