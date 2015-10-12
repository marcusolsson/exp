package main

import (
	"flag"
	"log"
)

var (
	defaultPort = "3001"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	var bindAddr, joinAddr string

	flag.StringVar(&bindAddr, "bind", "0.0.0.0:"+defaultPort, "")
	flag.StringVar(&joinAddr, "join", "", "")
	flag.Parse()

	srv := NewServer(bindAddr)
	srv.Start()

	if joinAddr != "" {
		srv.Join(joinAddr)
	}

	if err := srv.Listen(); err != nil {
		log.Fatal(err)
	}
}
