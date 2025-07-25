package main

import (
	"github.com/steadyfall/deris"
	"flag"
)

func main() {
	var srvPort uint64
	flag.Uint64Var(&srvPort, "port", 6969, "port to listen on")
	flag.Parse()
	
	deris.StartServer(srvPort, "snapshot.log")
}
