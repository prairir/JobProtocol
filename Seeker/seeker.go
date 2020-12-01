package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"

	globals "github.com/prairir/JobProtocol/Globals"
	jobs "github.com/prairir/JobProtocol/Jobs"
)

func main() {
	// set timeout and connection
	conn, err := net.Dial(globals.ConnType, fmt.Sprint(globals.ConnAddr, ":", globals.ConnPort))
	fatalErrorCheck(err)

	// state
	// 0 initial connection
	// 1 waiting for HELLOACK
	// 2 first JOB EQN
	// 3 accepted and waiting for second JOB EQN
	// 4 closed
	state := 0
	// send HELLO at initial connection
	fmt.Fprintln(conn, "HELLO")
	fmt.Println("sent HELLO")
	state = 1
	for {
		fmt.Println("waiting for creator...")
		result, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}
		cleanedResult := strings.TrimSpace(string(result))
		fmt.Println("received:", cleanedResult, "state:", state)

		if state == 1 && cleanedResult == "HELLOACK" {
			state = 2
		} else if state == 2 && cleanedResult == "JOB EQN" {
			conn.Write([]byte("ACPT JOB EQN\r\n"))
			fmt.Println("accept:", result)
			state = 3
		} else if state == 2 && cleanedResult == "JOB TCPFLOOD" {
			conn.Write([]byte("ACPT JOB TCPFLOOD\r\n"))
			fmt.Println("accept:", result)
			state = 3
		} else if state == 2 && cleanedResult == "JOB HOSTUP" {
			conn.Write([]byte("ACPT JOB HOSTUP\r\n"))
			fmt.Println("accept:", result)
			state = 3
		} else if state == 3 && cleanedResult[:7] == "JOB EQN" {
			fmt.Println("[", cleanedResult[8:], "]")
			expression, err := govaluate.NewEvaluableExpression(cleanedResult[8:])
			if err != nil {
				fmt.Println("job failed... bad input?", err.Error())
				conn.Write([]byte("JOB FAIL\r\n"))
			}
			expResult, err := expression.Evaluate(nil)
			if err != nil {
				fmt.Println("job failed... bad input?", err.Error())
				conn.Write([]byte("JOB FAIL\r\n"))
			} else {
				conn.Write([]byte("JOB SUCC " + fmt.Sprint(expResult) + "\r\n"))
			}
			state = 4
		} else if state == 3 && cleanedResult[:12] == "JOB TCPFLOOD" {
			// splits after JOB TCPFLOOD
			// eg JOBTCPFLOOD 123.321.543.345 14 -> ["123.321.543.345", "14"]
			splits := strings.Split(cleanedResult[:13], " ")
			port, _ := strconv.Atoi(splits[1])

			jobs.TCPFlood(splits[0], port)

			conn.Write([]byte("JOB SUCC \r\n"))
			state = 4
		} else if state == 3 && cleanedResult[:10] == "JOB HOSTUP" {
			hosts := strings.Split(cleanedResult[11:], " ")
			fmt.Println(hosts)
			var buf bytes.Buffer
			for _, host := range hosts {
				online, offline, err := jobs.HostUp(host, &buf)
				if err != nil {
					if err.Error() != "timeout" {
						conn.Write([]byte("JOB FAIL"))
						break
					}
				}
				conn.Write([]byte("JOB SUCC ONLINE "))
				for _, ip := range online {
					conn.Write([]byte(fmt.Sprint(ip, ", ")))
				}
				conn.Write([]byte("OFFLINE "))
				for _, ip := range offline {
					conn.Write([]byte(fmt.Sprint(ip, ", ")))
				}
				conn.Write([]byte("\r\n"))
				break
			}
			state = 4
		}
		if state == 4 {
			state = 2
			continue
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
