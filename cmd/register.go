package main

import (
	"flag"
	"log"

	"github.com/gwbeacon/gwbeacon/server/register"
)

func main() {
	var port int
	var timeBase int64
	flag.IntVar(&port, "port", 9999, "-port 9999")
	flag.Int64Var(&timeBase, "timebase", 0, "-timebase 1534061219")
	flag.Parse()
	server := register.NewServer(port, timeBase)

	err := server.Serve()
	log.Println(err)
}
