package main

import (
	"fmt"
	"strconv"

	"./server"
)

func main() {
	host := "127.0.0.1"
	port := 25565

	fmt.Println("Starting goserve...")

	server := server.CreateServer(host, port)
	server.Start()

	fmt.Println("goserve running on " + server.Host + ":" + strconv.Itoa(server.Port) + "...")
}
