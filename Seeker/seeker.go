package main

import (
	"bufio"
	"container/list"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	// the port that we use
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter port:")
	port, _ := reader.ReadString('\n')
	// add the address and get rid of the \n at the end of port
	addr := "127.0.0.1:" + port[:len(port)-1]
	fmt.Println(addr)

	// open port for tcp connection
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	checkError(err)

	// create a listener on that open port
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// couldnt figure out a channel version
	// its probably not memory safe but who knows
	queue := list.New()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handleConnection(conn, queue)

		daytime := time.Now().String()
		conn.Write([]byte(daytime)) // don't care about return value
		conn.Close()                // we're finished with this client
	}
}

func handleConnection(conn net.Conn, queue *list.List) {
	// state values
	// 0 = connection initalized
	// 1 connection established
	// 2 avaible
	// 3 ready for jobs
	// 4 closed
	state := 0
	for {

	}

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
