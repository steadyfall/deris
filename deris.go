package main

import (
	"bufio"
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		input := scanner.Text()
		fmt.Printf("got ping... [message='%s']\n", input)
		fmt.Fprintln(conn, "Hello World")
	}
}

func main() {
	lis, err := net.Listen("tcp", ":6969")
	if err != nil {
		panic(err)
	}
	fmt.Println("server started on port 6969")
	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}
