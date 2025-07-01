package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func handleError(c net.Conn, errInfo string) {
	errStr := fmt.Sprintf("ERROR: %s", errInfo)
	fmt.Println(errStr)
	fmt.Fprintln(c, errStr)
}

func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("got ping... [cmd='%s']\n", line)

		args := strings.Fields(line)
		if len(args) == 0 {
			handleError(conn, "empty command")
			continue
		}
		cmd := strings.ToUpper(args[0])

		switch cmd {
		case "GET":
			if len(args[1:]) != 1 {
				handleError(conn, "invalid number of arguments")
				continue
			}
			fmt.Fprintf(conn, "%s\n", args[1])
		case "SET":
			if len(args[1:]) != 2 {
				handleError(conn, "invalid number of arguments")
				continue
			}
			fmt.Fprintf(conn, "%s-%s\n", args[1], args[2])
		case "DEL":
			if len(args[1:]) != 1 {
				handleError(conn, "invalid number of arguments")
				continue
			}
			fmt.Fprintf(conn, "%s\n", args[1])
		}
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
