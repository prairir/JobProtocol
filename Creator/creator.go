package main

import (
	"bufio"
	"fmt"
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

	go jobSender(&mutex, &queue)
	for {
		conn, err := listener.Accept()
		// if error go through close connection process
		if err != nil {
			fmt.Println("Could not accept TCP connection! ", err.Error())
			os.Exit(1)
		}
		defer conn.Close()
		go handleHello(conn, &mutex, &queue)
	}
}
// state values
// 0 waiting for HELLO 
// 1 HELLO received, add element to queue 
// -- PROCESSING JOB --
// 2 JOB accepted/rejected
// 3 JOB result
func jobSender(mutex *sync.Mutex, queue *[]net.Conn) {
	state := 0
	var cleanedResult string
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter query: ")
		eqn, _ := reader.ReadString('\n')
		var conn net.Conn
		mutex.Lock()
		fmt.Println(*queue)
		if len(*queue) != 0 {
			conn = (*queue)[0]
			*queue = (*queue)[1:]
			fmt.Println(*queue)
		} else {
			fmt.Println("No jobs available, please try again later. ")
			mutex.Unlock()
			continue
		}
		mutex.Unlock()
		fmt.Println("conn:", conn, "state:", state, "eqn:", eqn)
		return

		if strings.Compare(cleanedResult, "JOB TIME") == 0 {
			daytime := time.Now().String()
			conn.Write([]byte("DONE TIME " + daytime))
		} else if strings.Compare(cleanedResult[:6], "JOB EQ") == 0 { // evaluate equation
			expression, _ := govaluate.NewEvaluableExpression(cleanedResult[6:])
			expResult, _ := expression.Evaluate(nil)
			conn.Write([]byte("DONE EQ " + expResult.(string)))
		}
	}
}

func handleHello(conn net.Conn, mutex *sync.Mutex, queue *[]net.Conn) {
	// event loop
	result, _ := bufio.NewReader(conn).ReadString('\n')
	cleanedResult := strings.TrimSpace(string(result))
	
	// initial message handling
	if strings.Compare(cleanedResult, "HELLO") == 0 {
		fmt.Fprintf(conn, "HELLOACK")
		mutex.Lock()
		*queue = append(*queue, conn)
		fmt.Println(*queue)
		mutex.Unlock()
	}
	return
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
