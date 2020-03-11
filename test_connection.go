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
type Bar struct {
	X int
}

type baz struct {
	Y int
}

type foo struct {
	Bar      // anonymous struct fields must also be exported types
	Baz  baz // non-anonymous structs can be unexported types
	A, B int
	Msg  string
	Map  map[string]int
}

func (f *foo) String() string {
	return fmt.Sprintf("{%d}, %v, %d, %d: %s | %v", f.X, f.Baz, f.A, f.B, f.Msg, f.Map)
}

func main() {
	addr := flag.String("addr", "localhost", "the destination address; only used as client mode")
	port := flag.String("port", "8000", "port to use")
	client := flag.Bool("client", false, "if the program is running as client (otherwise runs as server)")
	flag.Parse()

	f := foo{Bar{19}, baz{33}, 7, 42, "hello, world", map[string]int{"hello": 1, "hi": 2, "hey": 3}}

	if *client {
		conn, err := net.Dial("tcp", *addr+":"+*port)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		enc := gob.NewEncoder(conn)
		_ = enc.Encode(f) // look we're not programming well here
		fmt.Println("remote:", conn.RemoteAddr())
		fmt.Println("local: ", conn.LocalAddr())
	} else {
		listener, _ := net.Listen("tcp", ":"+*port)
		conn, _ := listener.Accept()
		defer conn.Close()
		dec := gob.NewDecoder(conn)
		f2 := new(foo)
		_ = dec.Decode(&f2)
		fmt.Println(f2)
		fmt.Println("remote:", conn.RemoteAddr())
		fmt.Println("local: ", conn.LocalAddr())
	}
}
