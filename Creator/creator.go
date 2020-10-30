package main

import (
	"fmt"
	"net"
	"os"
)

const (
	CONN_ADDR = "localhost"
	CONN_PORT = 7777
	CONN_TYPE = "tcp"
)

func main() {
	_, err := net.Listen(CONN_TYPE, fmt.Sprint(CONN_ADDR, ":", CONN_PORT))
	checkError(err)
	fmt.Println("listening on port ", CONN_PORT)
	for {

	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
