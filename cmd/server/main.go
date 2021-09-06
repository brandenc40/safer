package main

import (
	"log"

	"github.com/brandenc40/safer/internal/restserver"
)

func main() {
	r := restserver.NewServer()
	log.Fatal(r.StartServer(":8080"))
}
