package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"github.com/Knetic/govaluate"
)

const (
	CONN_ADDR = "localhost"
	CONN_PORT = 7777
	CONN_TYPE = "tcp"
)

func main() {
	// set timeout and connection
	conn, err := net.Dial(CONN_TYPE, fmt.Sprint(CONN_ADDR, ":", CONN_PORT))
	fatalErrorCheck(err)

	// state
	// 0 initial connection
	// 1 waiting for HELLOACK
	// 2 first JOB EQN
	// 3 accepted and waiting for second JOB EQN
	// 4 closed
	state := 0
	for {
		fmt.Println("test1")
		result, err := ioutil.ReadAll(conn)
		fmt.Println("test2")
		// clean the result and avoid error
		if err != nil {
			continue
		} 
		cleanedResult := strings.TrimSpace(string(result))

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

func fatalErrorCheck(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
