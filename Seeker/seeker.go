package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

type IdvSession struct {
	id         int64
	state      int
	connection net.Conn
}

func main() {
	// read port from console
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
		// if error go through close connection process
		if err != nil {
			closeConnection(conn)
		}

		go handleConnection(conn, queue)

		// just for testing purposes
		daytime := time.Now().String()
		conn.Write([]byte(daytime)) // don't care about return value
		conn.Close()                // we're finished with this client
	}
}

func handleConnection(conn net.Conn, queue *list.List) {
	// state values
	// 0 = connection initalized
	// 1 connection established
	// 2 full
	// 3 ready for jobs
	// 4 closed
	state := 0
	currSession := IdvSession{
		// session ID is gonna be the time in nano seconds
		// I thought bout writing a mathematically unique ID
		// but the chances are crazy low soooooooo I aint doing it
		time.Now().UnixNano(),
		state,
		conn,
	}
	queue.PushBack(currSession)
	// event loop
	for {
		result, err := ioutil.ReadAll(conn)
		if err != nil {
			closeConnection(conn)
		}

		cleanedResult := strings.TrimSpace(string(result))

		// set state based on position
		var position int = getPosition(&currSession, queue)
		if position == 0 {
			state = 3
		} else {
			state = 2
		}

		// initial message handling
		if state == 0 && strings.Compare(cleanedResult, "HELLO") == 0 {
			state = 1
			conn.Write([]byte("HELLOACK"))
		} else if state == 2 { // options for when its not at front of queue
			if strings.Compare(cleanedResult, "AVL") == 0 {
				conn.Write([]byte("FULL"))
			}
		} else if state == 3 { // options for when it IS at the front of the queue
			if strings.Compare(cleanedResult, "AVL") == 0 {
				conn.Write([]byte("AVLACK"))
			} else if strings.Compare(cleanedRestult, "JOB TIME") == 0 {
				daytime := time.Now().String()
				conn.Write([]byte(daytime))
			} else if strings.Compare(cleanedRestult[:6], "JOB EQ") == 0 {
				// equation checker stuff here
			}
		}
		daytime := time.Now().String()
		conn.Write([]byte(daytime))
	}

}

func getPosition(currSession *IdvSession, queue *list.List) int {

}

func closeConnection(conn net.Conn) {
	conn.Write([]byte("BYE"))
	conn.Close()
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
