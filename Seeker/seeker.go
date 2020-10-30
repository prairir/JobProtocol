package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"

	"github.com/Knetic/govaluate"
)

func main() {
	// read port from console
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter address to connect to: ")
	addr, _ := reader.ReadString('\n')

	fmt.Print("Enter port to connect to: ")
	port, _ := reader.ReadString('\n')
	// add the address and get rid of the \n at the end
	fullAddr := addr[:len(addr)-1] + ":" + port[:len(port)-1]
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
				break
			}
			expResult, err := expression.Evaluate(nil)
			if err != nil {
				conn.Write([]byte("JOB FAIL"))
				break
			} else {
				conn.Write([]byte("JOB SUCC " + expResult.(string)))
			}
			state = 4
		}
		if state == 4 {
			break
		}
	}
	conn.Close()
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
