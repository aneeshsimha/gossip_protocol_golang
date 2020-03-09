package main

//76.21.104.153

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net"
)

// ONLY exported fields are transmitted using gob
type foo struct {
	A, B int
	Msg  string
	Map  map[string]int
}

func (f *foo) String() string {
	return fmt.Sprintf("%d, %d: %s | %v", f.A, f.B, f.Msg, f.Map)
}

func main() {
	addr := flag.String("addr", "localhost", "the destination address; only used as client mode")
	port := flag.String("port", "8000", "port to use")
	client := flag.Bool("client", false, "if the program is running as client (otherwise runs as server)")
	flag.Parse()

	f := foo{7, 42, "hello, world", map[string]int{"hello": 1, "hi": 2, "hey": 3}}

	if *client {
		conn, err := net.Dial("tcp", *addr+":"+*port)
		if err != nil {
			log.Fatal(err)
		}
		enc := gob.NewEncoder(conn)
		enc.Encode(f)
	} else {
		listener, _ := net.Listen("tcp", ":"+*port)
		conn, _ := listener.Accept()
		dec := gob.NewDecoder(conn)
		f2 := new(foo)
		dec.Decode(&f2)
		fmt.Println(f2)
	}
}
