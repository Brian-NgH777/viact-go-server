package main

import (
	"fmt"
	s "go/server/services"
)

const (
	port = 8700
)

func main() {

	fmt.Println("http server run port:", port)
	server := s.New()
	server.FastHttp(port)
}


