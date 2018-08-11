package main

import (
	"flag"
	"log"

	"github.com/gwbeacon/gwbeacon/connector"
)

func main() {
	var port int
	var registerAddr = ""
	flag.IntVar(&port, "port", 8888, "-port 8888")
	flag.StringVar(&registerAddr, "register", "localhost:9999", "-register localhost:9999")
	flag.Parse()
	server := connector.NewServer(int32(port), registerAddr)
	err := server.Serve()
	log.Println(err)
}
