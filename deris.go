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

func handleConnection(conn net.Conn, kv map[string]string) {
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
		// cmd: GET <key>
		case "GET":
			if len(args[1:]) != 1 {
				handleError(conn, "invalid number of arguments")
				continue
			}
			val, present := kv[args[1]]
			if !present {
				handleError(conn, fmt.Sprintf("key '%s' not present", args[1]))
				continue
			}
			fmt.Fprintf(conn, "%s\n", val)

		// cmd: SET <key> <value>
		case "SET":
			if len(args[1:]) != 2 {
				handleError(conn, "invalid number of arguments")
				continue
			}
			// TODO: Make SET operation immutable
			key, val := args[1], args[2]
			kv[key] = val

		// cmd: DEL <key>
		case "DEL":
			if len(args[1:]) != 1 {
				handleError(conn, "invalid number of arguments")
				continue
			}
			_, present := kv[args[1]]
			if !present {
				handleError(conn, fmt.Sprintf("key '%s' not present", args[1]))
				continue
			}
			delete(kv, args[1])
		}
	}
}

func main() {
	db := make(map[string]string)
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
		go handleConnection(conn, db)
	}
}
