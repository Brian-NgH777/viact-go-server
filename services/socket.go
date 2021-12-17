package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

const (
	portSocket = 8888
	hostSocket = "localhost"
)

type M struct {
	Key      string `json:"key"`
	Datetime string `json:"datetime"`
	GPS      string `json:"gps"`
}

func SocketServer() {
	serviceSocket := fmt.Sprintf("%s:%d", hostSocket, portSocket)
	fmt.Println("socket server run service: %s", serviceSocket)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", serviceSocket)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	defer listener.Close()

	r := RConnect()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn, r)
	}

}

func handleClient(conn net.Conn, redis *Redis) {
	// close connection on exit
	defer conn.Close()

	var buf [512]byte
	for {
		// read upto 512 bytes
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}

		msg := &M{}
		if err := json.Unmarshal(buf[0:n], msg); err != nil {
			return
		}

		_, err = redis.HSet(msg.Key, msg.GPS, fmt.Sprintf("%s-%s", msg.GPS, msg.Datetime))
		if err != nil {
			continue
		}

		data, err := redis.HGet(msg.Key, msg.GPS)
		if err != nil {
			continue
		}
		fmt.Println("HGet", data)
		reply := fmt.Sprintf("Done Data: %s", data)

		// write the n bytes read
		_, err2 := conn.Write([]byte(reply))
		if err2 != nil {
			return
		}
	}
}

func SocketClient() {
	service := "210.3.19.86:5566"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	_, err = conn.Write([]byte("test"))
	checkError(err)
	result, err := ioutil.ReadAll(conn)
	checkError(err)
	fmt.Println(string(result))
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
