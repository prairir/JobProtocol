package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"sync"
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

	// the mutex makes it memory safe
	var mutex sync.Mutex
	queue := list.New()

	for {
		conn, err := listener.Accept()
		// if error go through close connection process
		if err != nil {
			closeConnection(conn)
		}

		go handleConnection(conn, &mutex, queue)
	}
}

func handleConnection(conn net.Conn, mutex *sync.Mutex, queue *list.List) {
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
	mutex.Lock()
	queue.PushBack(currSession)
	mutex.Unlock()
	// event loop
	for {
		result, err := ioutil.ReadAll(conn)
		if err != nil {
			closeConnection(conn)
			state = 4
		}

		cleanedResult := strings.TrimSpace(string(result))

		// set state based on position
		position, entry := getPosition(&currSession, mutex, queue)
		if position == 0 {
			state = 3
		} else {
			state = 2
		}

		if state == 4 {
			mutex.Lock()
			queue.Remove(entry)
			mutex.Unlock()
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
			} else if strings.Compare(cleanedResult, "JOB TIME") == 0 {
				daytime := time.Now().String()
				conn.Write([]byte("DONE TIME " + daytime))
			} else if strings.Compare(cleanedResult[:6], "JOB EQ") == 0 {
				// equation checker stuff here
			}
		} else if strings.Compare(cleanedResult, "BYE") == 0 {
			state = 4
			continue
		}

		if state == 4 {
			break
		}
	}

}

func getPosition(currSession *IdvSession, mutex *sync.Mutex, queue *list.List) (int, *list.Element) {
	mutex.Lock()
	// Iterate through list and print its contents.
	// ya ya ya i know it could be a binary search but i dont have enough time to write that
	position := 0
	var e *list.Element
	for e, index := queue.Front(), 0; e != nil; e, index = e.Next(), index+1 {
		fmt.Println(e.Value)
		// checking the value of the current session to the session in queue
		if currSession.id == e.Value.(*IdvSession).id {
			position = index
			break
		}
	}
	mutex.Unlock()
	return position, e
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
