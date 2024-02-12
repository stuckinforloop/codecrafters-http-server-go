package main

import (
	"flag"
	"log"

	"github.com/codecrafters-io/http-server-starter-go/internal/server"
)

func main() {
	dir := flag.String("directory", "./", "directory containing files")
	flag.Parse()

	s := server.New(":4221", *dir, "tcp")
	if err := s.Start(); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
