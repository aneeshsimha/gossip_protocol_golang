package main

import (
	"coen317/gossip/gossip"
	"log"
)

func main() {
	gc := &gossip.Client{}

	log.Fatal(gc.Run())
}
