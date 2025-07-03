package deris

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func loadSnapshot(kv map[string]string, logFile string) {
	fmt.Printf("loading data from %s :)\n", logFile)
	snapshot, _ := os.Open(logFile)
	scanner := bufio.NewScanner(snapshot)
	for scanner.Scan() {
		line := scanner.Text()
		args := strings.Fields(line)
		cmd := strings.ToUpper(args[0])
		switch cmd {
		case "SET":
			key, val := args[1], args[2]
			kv[key] = val
		case "DEL":
			delete(kv, args[1])
		}
	}
}

func handleError(w io.Writer, errInfo string) {
	errStr := fmt.Sprintf("ERROR: %s", errInfo)
	fmt.Println(errStr)
	fmt.Fprintln(w, errStr)
}

func handleConnection(conn net.Conn, kv map[string]string, logFile io.Writer) {
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
			fmt.Fprintf(logFile, fmt.Sprintf("%s %s %s\n", args[0], args[1], args[2]))

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
			fmt.Fprintf(logFile, fmt.Sprintf("%s %s\n", args[0], args[1]))

		case "EXIT":
			fmt.Fprintln(conn, "-- Quitting session...")
			conn.Close()
		}
	}
}

func StartServer(port uint64) {
	db := make(map[string]string)

	if _, err := os.Stat("snapshot.log"); err == nil {
		loadSnapshot(db, "snapshot.log")
	} else if os.IsNotExist(err) {
		fmt.Println("snapshot not available :|")
	}
	snapshot, _ := os.OpenFile("snapshot.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	fmt.Println("server started on port 6969")

	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn, db, snapshot)
	}
}

func main() {
	var srvPort uint64 = 6969
	StartServer(srvPort)
}
