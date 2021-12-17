package main

import (
	"fmt"
	s "go/server/services"
)

const (
	port = 8700
	host = "0.0.0.0"
)

func main() {

	fmt.Println("http server run port:", port)
	server := s.New()
	server.FastHttp(host, port)
}


