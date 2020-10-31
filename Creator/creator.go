package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Knetic/govaluate"
)

/*
type IdvSession struct {
	id         int64
	state      int
	connection net.Conn
}
*/

const (
	CONN_ADDR = "localhost"
	CONN_PORT = 7777
	CONN_TYPE = "tcp"
)

func main() {
	// create a listener on that open port
	listener, err := net.Listen(CONN_TYPE, fmt.Sprint(CONN_ADDR, ":", CONN_PORT))
	fatalErrorCheck(err)
	defer listener.Close()
	fmt.Println("listening to", CONN_ADDR, "at port", CONN_PORT)

	
	var mutex sync.Mutex
	var queue []net.Conn
	for {
		conn, err := listener.Accept()
		// if error go through close connection process
		if err != nil {
			fmt.Println("Could not accept TCP connection! ", err.Error())
			os.Exit(1)
		}
		defer conn.Close()
		go handleHello(conn, &mutex, queue)
	}
}
// state values
// 0 waiting for HELLO 
// 1 HELLO received, add element to queue 
// -- PROCESSING JOB --
// 2 JOB accepted/rejected
// 3 JOB result

func handleHello(conn net.Conn, mutex *sync.Mutex, queue []net.Conn) {
	state := 0
	// event loop
	for {
		result, err := ioutil.ReadAll(conn)
		if err != nil {
			fmt.Println("could not read from client: ", err.Error())
		}
		cleanedResult := strings.TrimSpace(string(result))
		fmt.Println(cleanedResult)
		return

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
		}
	}
}

/*
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
*/

func fatalErrorCheck(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
