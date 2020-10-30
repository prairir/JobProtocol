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

	"github.com/Knetic/govaluate"
)

func main() {
	// read port from console
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter address to connect to:")
	addr, _ := reader.ReadString('\n')

	fmt.Print("Enter port to connect to:")
	port, _ := reader.ReadString('\n')
	// add the address and get rid of the \n at the end
	fullAddr := addr[:len(addr)-1] + port[:len(port)-1]
	fmt.Println(addr)

	// set timeout and connection
	timeout, err := time.ParseDuration("5s")
	conn, err := net.DialTimeout("tcp", fullAddr, timeout)
	checkError(err)

	// state
	// 0 initial connection
	// 1 waiting for HELLOACK
	// 2 first JOB EQN
	// 3 accepted and waiting for second JOB EQN
	// 4 closed
	state := 0
	for {
		result, err := ioutil.ReadAll(conn)
		// clean the result and avoid error
		var cleanedResult string
		if err != nil {
			continue
		} else {
			cleanedResult = strings.TrimSpace(string(result))
		}

		// send HELLO at initial connection
		if state == 0 {
			conn.Write([]byte("HELLO"))
		} else if state == 1 && cleanedResult == "HELLOACK" {
			state = 2
		} else if state == 2 && cleanedResult == "JOB EQN" {
			conn.Write([]byte("ACPT JOB EQN"))
			state = 3
		} else if state == 3 && cleanedResult[:7] == "JOB EQN" {
			expression, err := govaluate.NewEvaluableExpression(cleanedResult[6:])
			if err != nil {
				conn.Write([]byte("JOB FAIL"))
			}
			expResult, err := expression.Evaluate(nil)
			if err != nil {
				conn.Write([]byte("JOB FAIL"))
			}
			conn.Write([]byte("DONE EQ " + expResult.(string)))
		}
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
			} else if strings.Compare(cleanedResult[:6], "JOB EQ") == 0 { // evaluate equation
				expression, _ := govaluate.NewEvaluableExpression(cleanedResult[6:])
				expResult, _ := expression.Evaluate(nil)
				conn.Write([]byte("DONE EQ " + expResult.(string)))
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
