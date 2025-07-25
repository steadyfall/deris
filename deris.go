package deris

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
)

func loadSnapshot(kv map[string]string, logFile string) error {
	fmt.Printf("loading data from %s :)\n", logFile)
	snapshot, err := os.OpenFile(logFile, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error opening AOF file for R/W: %w", err)
	}
	defer snapshot.Close()

	var expectedChecksum string
	scanner := bufio.NewScanner(snapshot)
	hasher := sha256.New()
	contentSize := int64(0)
	for scanner.Scan() {
		line := scanner.Text()
		args := strings.Fields(line)
		cmd := strings.ToUpper(args[0])
		switch cmd {
		case "SET":
			if len(args[1:]) != 2 {
				return fmt.Errorf("%s | %s", line, "invalid number of arguments")
			}
			n, err := hasher.Write([]byte(line + "\n"))
			if err != nil {
				return fmt.Errorf("failed to hash content: %w", err)
			}
			contentSize += int64(n)
			key, val := args[1], args[2]
			kv[key] = val
		case "DEL":
			if len(args[1:]) != 1 {
				return fmt.Errorf("%s | %s", line, "invalid number of arguments")
			}
			n, err := hasher.Write([]byte(line + "\n"))
			if err != nil {
				return fmt.Errorf("failed to hash content: %w", err)
			}
			contentSize += int64(n)
			delete(kv, args[1])
		case "SHA256:":
			if len(args[1:]) != 1 {
				return fmt.Errorf("%s | %s", line, "invalid number of arguments")
			}
			expectedChecksum = args[1]
			if len(expectedChecksum) != 64 {
				return fmt.Errorf("sha256 integrity not valid")
			}
		}
	}

	actualHash := hasher.Sum(nil)
	actualChecksum := fmt.Sprintf("%x", actualHash)
	if expectedChecksum != actualChecksum {
		return fmt.Errorf("sha256 integrity not satisfied, snapshot corrupted")
	}
	hasher.Reset()

	err = snapshot.Truncate(contentSize)
	if err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}
	err = snapshot.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync file after truncation: %w", err)
	}

	return nil
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
			cmdLine := fmt.Sprintf("%s %s %s", cmd, args[1], args[2])
			key, val := args[1], args[2]
			kv[key] = val
			fmt.Fprintln(logFile, cmdLine)

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
			cmdLine := fmt.Sprintf("%s %s", cmd, args[1])
			delete(kv, args[1])
			fmt.Fprintln(logFile, cmdLine)

		case "EXIT":
			fmt.Fprintln(conn, "-- Quitting session...")
			conn.Close()
		}
	}
}

func handleShutdown(c <-chan os.Signal, l net.Listener, logFile string, openLogFile io.Writer) {
	<-c
	fmt.Println("exiting server...")
	l.Close()

	fileContent, err := os.ReadFile(logFile)
	if err != nil {
		fmt.Println("ERROR: snapshot corrupted")
		os.Exit(1)
	}
	checksum := sha256.Sum256(fileContent)
	fmt.Fprintf(openLogFile, "SHA256: %x\n", checksum)

	os.Exit(0)
}

func StartServer(port uint64, snapshotFile string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	db := make(map[string]string)

	if _, err := os.Stat(snapshotFile); err == nil {
		err = loadSnapshot(db, snapshotFile)
		if err != nil {
			fmt.Printf("CRITICAL: %s\n", err.Error())
			os.Exit(2)
		}
	} else if os.IsNotExist(err) {
		fmt.Println("snapshot not available :|")
	}
	snapshot, _ := os.OpenFile(snapshotFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer snapshot.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	fmt.Println("server started on port 6969")

	go handleShutdown(c, lis, snapshotFile, snapshot)

	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn, db, snapshot)
	}
}
