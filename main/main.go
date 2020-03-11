package main

import (
	"flag"
	"fmt"
	"github.com/aneeshsimha/gossip_protocol_golang/gossip"
	"log"
	"math/rand"
	"time"
)

func main() {
	addr := flag.String("addr", "", "the known address of a node to join the network through")
	alivePort := flag.String("alive", "8000", "port for keep alives")
	msgPort := flag.String("msgPort", "8001", "port for message passing")
	loops := flag.Uint64("loops", 10, "number of seconds to loop for")
	flag.Parse()
	sendTime := rand.Uint64() % (*loops - 1) + 1
	log.Printf("looping for %d seconds, sending a message at %d seconds", *loops, sendTime)

	gc := gossip.New(
		3,
		3,
		*alivePort,
		*msgPort,
		200*time.Millisecond,
		200*time.Millisecond,
	)
	gc.Run(*addr)

	time.Sleep(time.Duration(sendTime) * time.Second)
	gc.Send(fmt.Sprintf("a message @ %v", sendTime))
	time.Sleep(time.Duration(*loops - sendTime) * time.Second)

	gc.Shutdown()
	time.Sleep(time.Second)

	for _, e := range gc.Nodes() {
		fmt.Println(e)
	}

	for _, e := range gc.Messages() {
		fmt.Println(e)
	}
}
