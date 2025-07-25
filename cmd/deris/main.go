package main

import (
	"github.com/steadyfall/deris"
)

func main() {
	var srvPort uint64 = 6969
	deris.StartServer(srvPort, "snapshot.log")
}
