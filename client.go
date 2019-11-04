package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const host, port = "127.0.0.1", "1337"

func main() {
	conn, err := net.Dial("tcp", host+":"+port)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer conn.Close()

	go ProcessServerMessages(conn)
	ProcessConsoleInput(conn)
}

// ProcessServerMessages handles server output
func ProcessServerMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for true {
		message, _ := reader.ReadString('\n')
		message = strings.TrimSuffix(message, "\n")

		if message != "" {
			fmt.Println(message)
		}
	}
}

// ProcessConsoleInput handles console input
func ProcessConsoleInput(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)

	for true {
		message, _ := reader.ReadString('\n')
		message = strings.TrimSuffix(message, "\n")

		if message != "" {
			fmt.Fprintln(conn, message)
		}
	}
}
