package main

import (
	"flag"
	"fmt"
	"github.com/aneeshsimha/gossip_protocol_golang/gossip"
	"log"
	"math/rand"
	"time"
)

const (
	Black        = "\033[0;30m"
	Red          = "\033[0;31m"
	Green        = "\033[0;32m"
	Brown_Orange = "\033[0;33m"
	Blue         = "\033[0;34m"
	Purple       = "\033[0;35m"
	Cyan         = "\033[0;36m"
	Light_Gray   = "\033[0;37m"
	Dark_Gray    = "\033[1;30m"
	Light_Red    = "\033[1;31m"
	Light_Green  = "\033[1;32m"
	Yellow       = "\033[1;33m"
	Light_Blue   = "\033[1;34m"
	Light_Purple = "\033[1;35m"
	Light_Cyan   = "\033[1;36m"
	White        = "\033[1;37m"
	CLEAR        = "\033[0m"
)

var COLORS = []string{
	Red,
	Green,
	Blue,
	Purple,
	Cyan,
	Yellow,
}

func main() {
	addr := flag.String("addr", "", "the known address of a node to join the network through")
	alivePort := flag.String("alive", "8000", "port for keep alives")
	msgPort := flag.String("msgPort", "8001", "port for message passing")
	loops := flag.Uint64("loops", 10, "number of seconds to loop for")
	flag.Parse()

	gc := gossip.New(
		3,
		3,
		*alivePort,
		*msgPort,
		50*time.Millisecond,
		50*time.Millisecond,
	)
	gc.Run(*addr)

	rand.Seed(time.Now().UnixNano())
	sendTime := rand.Uint64()%(*loops-1) + 1
	log.Printf("looping for %d seconds, sending a message at %d seconds", *loops, sendTime)

	time.Sleep(time.Duration(sendTime) * time.Second)
	//gc.Send(fmt.Sprintf("a message @ %v", sendTime))
	gc.Send(fmt.Sprintf("%sa colorful message @ %v%s", COLORS[rand.Int() % len(COLORS)], sendTime, CLEAR))
	time.Sleep(time.Duration(*loops-sendTime) * time.Second)

	gc.Shutdown()
	time.Sleep(time.Second)

	for _, e := range gc.Nodes() {
		fmt.Println(e)
	}

	for _, e := range gc.Messages() {
		fmt.Println(e)
	}
	log.Printf("looped for %d seconds, sent own message at %d seconds", *loops, sendTime)
}
