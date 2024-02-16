package main

import (
	"context"
	"flag"
	"log"

	"github.com/codecrafters-io/http-server-starter-go/internal/server"
)

var (
	directory string
)

const (
	addr    string = ":4221"
	network string = "tcp"
)

func init() {
	dirFlag := flag.String("directory", "./", "directory containing files")
	flag.Parse()

	directory = *dirFlag
}

func main() {
	// TODO: use signal.NotifyContext to gracefully shutdown the server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := server.New(addr, directory, network)

	log.Printf("binding port %s for tcp connections", addr)
	if err := srv.Start(ctx); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
