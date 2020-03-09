package main

import (
	"coen317/gossip/gossip"
	"flag"
	"time"
)

func main() {
	addr := flag.String("addr", "localhost", "the known address of a node to join the network through")
	alivePort := flag.String("alive", "8000", "port for keep alives")
	msgPort := flag.String("msgPort", "8001", "port for message passing")
	flag.Parse()

	gc := gossip.New(
		5,
		5,
		*alivePort,
		*msgPort,
		100*time.Millisecond,
		100*time.Millisecond,
	)
	gc.Run(addr)
}
