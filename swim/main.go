package main

import (
	"flag"
	"log"
	"os"
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

	logger := log.New(os.Stdout, "swim: ", 0)

	srv := NewServer(bindAddr, logger)
	if err := srv.Start(); err != nil {
		logger.Fatal("unable to start server")
	}

	if joinAddr != "" {
		if err := srv.Join(joinAddr); err != nil {
			logger.Fatalf("unable to join %s", joinAddr)
		}
	}

	logger.Println("listening on", bindAddr)

	if err := srv.Listen(); err != nil {
		logger.Fatal(err)
	}
}
