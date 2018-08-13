package main

import (
	"flag"
	"log"
)

func main() {
	var port int
	var registerAddr = ""
	flag.IntVar(&port, "port", 6666, "-port 6666")
	flag.StringVar(&registerAddr, "register", "localhost:9999", "-register localhost:9999")
	flag.Parse()
	server := NewServer(int32(port), registerAddr)
	err := server.Serve()
	log.Println(err)
}
