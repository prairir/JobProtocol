package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	globals "github.com/prairir/JobProtocol/Globals"
)

func RunCreator(queueTR chan int, queueRV chan []net.Conn) {
	fmt.Println(globals.GetJobNames())
	// create a listener on that open port
	listener, err := net.Listen(globals.ConnType, fmt.Sprint(globals.ConnAddr, ":", globals.ConnPort))
	globals.FatalErrorCheck(err)
	defer listener.Close()
	fmt.Println("listening to", globals.ConnAddr, "at port", globals.ConnPort)

	c := Creator{}
	go c.cmd()
	for {
		// if error go through close connection process
		fmt.Println("Waiting for Job Seeker...")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Could not accept TCP connection!", err.Error())
			os.Exit(1)
		}
		defer conn.Close()
		go c.handleHello(conn)

		// return the queue
		select {
		case <-queueTR:
			queueRV <- c.queue
		default:
			break
		}
	}
}

// Creator type contains a queue of seekers to use as connections, and a mutex to access the queue.
type Creator struct {
	mutex sync.Mutex
	queue []net.Conn
}

// cmd for controlling the creator
//
// state values
// 0 waiting for HELLO
// 1 HELLO received, add element to queue
// -- PROCESSING JOB --
// 2 JOB accepted/rejected
// 3 JOB result
func (c *Creator) cmd() {
	var query string
	var header string
	var err error
	isNewQuery := true
	for {
		if isNewQuery {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter query: ")
			query, _ = reader.ReadString('\n')
			query = strings.TrimSpace(query)
			header, err = globals.GetHeader(query)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
		isNewQuery = true // variable reset
		var conn net.Conn
		c.mutex.Lock()
		fmt.Println(c.queue)
		if len(c.queue) != 0 {
			conn = (c.queue)[0]
			c.queue = (c.queue)[1:]
			fmt.Println(c.queue)
		} else {
			fmt.Println("No jobs available, please try again later. ")
			c.mutex.Unlock()
			continue
		}
		c.mutex.Unlock()
		fmt.Println(header)
		// job starts
		fmt.Fprintln(conn, header)
		accept, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Connection could not be established, try again with new seeker. ")
			conn.Close()
			isNewQuery = false
			continue
		}
		if len(accept) > 4 && accept[:4] != "ACPT" {
			fmt.Println("message not accepted. Trying with next connection. ")
			isNewQuery = false // use the same query on a new connection
			c.mutex.Lock()
			c.queue = append(c.queue, conn)
			c.mutex.Unlock()
			continue
		}

		// do the full query since the job was accepted
		fmt.Fprintln(conn, query)
		response, err := bufio.NewReader(conn).ReadString('\n')
		//if len(response) >= 10 && response[:8] == "JOB SUCC" {
		respHeader, err := globals.GetHeader(response)
		fmt.Println("response: [", respHeader, "]")
		if respHeader == "JOB SUCC" {
			fmt.Println("job done! result: ")
			fmt.Println(response[len(respHeader)+1:])
			c.mutex.Lock()
			c.queue = append(c.queue, conn)
			c.mutex.Unlock()
		} else {
			fmt.Println("job failed! trying with a new connection")
			conn.Close()
			isNewQuery = false
			continue
		}
	}
}

func (c *Creator) handleHello(conn net.Conn) {
	// event loop
	result, _ := bufio.NewReader(conn).ReadString('\n')
	cleanedResult := strings.TrimSpace(string(result))

	// initial message handling
	if strings.Compare(cleanedResult, "HELLO") == 0 {
		c.mutex.Lock()
		c.queue = append(c.queue, conn)
		fmt.Println(c.queue)
		c.mutex.Unlock()
		fmt.Fprintln(conn, "HELLOACK")
		fmt.Println("sent the HELLOACK")
	}
	return
}
