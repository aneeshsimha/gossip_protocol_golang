package main

import (
	"flag"
	"fmt"
	"github.com/aneeshsimha/gossip_protocol_golang/gossip"
	"time"
)

func main() {
	addr := flag.String("addr", "", "the known address of a node to join the network through")
	alivePort := flag.String("alive", "8000", "port for keep alives")
	msgPort := flag.String("msgPort", "8001", "port for message passing")
	flag.Parse()

	gc := gossip.New(
		3,
		3,
		*alivePort,
		*msgPort,
		1000*time.Millisecond,
		1000*time.Millisecond,
	)
	gc.Run(*addr)
	time.Sleep(10 * time.Second)

	gc.Shutdown()
	time.Sleep(time.Second)

	for _, e := range gc.Nodes() {
		fmt.Println(e)
	}
}
