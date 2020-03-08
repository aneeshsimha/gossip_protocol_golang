package main

import (
	"coen317/gossip/gossip"
	"log"
)

func main() {
	gc := &gossip.Client{}
	gc.Run()

	log.Fatal(gc.Run())
}
