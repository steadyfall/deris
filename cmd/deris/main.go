package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/steadyfall/deris"
)

func verifyDirsExist(path string) error {
	if path == "" {
		return fmt.Errorf("path is empty")
	}

	parent := filepath.Dir(filepath.Clean(path))

	info, err := os.Stat(parent)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", parent)
		}
		return fmt.Errorf("error accessing directory: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", parent)
	}

	return nil
}

func main() {
	var srvPort uint64
	flag.Uint64Var(&srvPort, "port", 6969, "port to listen on")

	var snapshotFile string
	flag.StringVar(&snapshotFile, "snapshot", "snapshot.log", "path of snapshot file")
	flag.Parse()

	// check to ensure that given directory exists
	if err := verifyDirsExist(snapshotFile); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid path: %v\n", err)
		os.Exit(1)
	}

	deris.StartServer(srvPort, snapshotFile)
}
