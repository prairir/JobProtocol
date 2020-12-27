package creator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	globals "github.com/prairir/JobProtocol/Globals"
)

//RunCreator runs the creator
func RunCreator(jobInput chan string, jobResult chan map[string]string, getQueue chan []net.Conn) {
	fmt.Println(globals.GetJobNames())
	// create a listener on that open port
	listener, err := net.Listen(globals.ConnType, fmt.Sprint(globals.ConnAddr, ":", globals.ConnPort))
	globals.FatalErrorCheck(err)
	defer listener.Close()
	fmt.Println("listening to", globals.ConnAddr, "at port", globals.ConnPort)

	c := Creator{}
	go c.cmd(jobInput, jobResult)
	go func() {
		for {
			getQueue <- c.queue
		}
	}()
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

	}
}

// Creator type contains a queue of seekers to use as connections, and a mutex to access the queue.
type Creator struct {
	mutex     sync.Mutex
	queue     []net.Conn
	MACbuffer map[string][]string
	custom    string // for custom creator commanders

}

func (c *Creator) firstPart(
	jobInput chan string,
	jobResult chan map[string]string,
	isNewQuery *bool,
	query *string,
	header *string,
	args *string,
) net.Conn {
	fmt.Println("ready for channel input...")
	var err error
	if *isNewQuery {
		// read input from the job channel
		*query = <-jobInput
		*query = strings.TrimSpace(*query)
		*header, err = globals.GetHeader(*query)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		*header = strings.ToLower(*header)
		*args = strings.ToLower(strings.TrimSpace((*query)[len(*header):]))
	}
	*isNewQuery = true // variable reset
	var conn net.Conn
	c.mutex.Lock()
	fmt.Println("queue:", c.queue)
	// getting top value of queue
	if len(c.queue) != 0 {
		conn = (c.queue)[0]
		c.queue = (c.queue)[1:]
		fmt.Println(c.queue)
	} else {
		fmt.Println("No jobs available, please try again later. ")
		c.mutex.Unlock()
		return nil
	}
	c.mutex.Unlock()
	fmt.Println("query:", *query)
	fmt.Println("header:", *header)
	fmt.Println("args:", *args)

	// job starts
	// sends header first
	if *query == "JOB GETMAC SPYSEEKER" {
		// re-synchronization
		*query = "JOB GETMAC"
		fmt.Fprintln(conn, *query)
	} else if *header == "job spyseeker" {
		fmt.Fprintln(conn, *header)
		*header = "job spy"
	}
	fmt.Fprintln(conn, *header)
	accept, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Connection could not be established, try again with new seeker. ")
		conn.Close()
		*isNewQuery = false
		return nil
	}
	if len(accept) > 4 && accept[:4] != "ACPT" {
		fmt.Println("message not accepted. Trying with next connection. ")
		*isNewQuery = false // use the same query on a new connection
		c.mutex.Lock()
		c.queue = append(c.queue, conn)
		c.mutex.Unlock()
		return nil
	}
	return conn
}

// cmd for controlling the creator
//
// state values
// 0 waiting for HELLO
// 1 HELLO received, add element to queue
// -- PROCESSING JOB --
// 2 JOB accepted/rejected
// 3 JOB result
func (c *Creator) cmd(jobInput chan string, jobResult chan map[string]string) {
	var query string
	var header string
	var args string
	isNewQuery := true
	for {
		conn := c.firstPart(jobInput, jobResult, &isNewQuery, &query, &header, &args)
		if conn == nil {
			continue
		}
		if header == "job spy" && args == "seeker" {
			jobInput2 := make(chan string, 1)
			jobResult2 := make(chan map[string]string)
			c.mutex.Lock()
			c.queue = append(c.queue, conn)
			c.mutex.Unlock()
			fmt.Println("RUNNING CMD...")
			go c.cmd(jobInput2, jobResult2)
			jobInput2 <- "JOB GETMAC SPYSEEKER"
			res := <-jobResult2
			fmt.Println("res:", res)
			var macs []string
			resStr := res[conn.RemoteAddr().String()]
			resStr = resStr[len("JOB GETMAC => "):]
			fmt.Println("resStr:", resStr)
			err := json.Unmarshal([]byte(resStr), &macs)
			if err != nil {
				fmt.Println(err)
				return
			}
			if c.MACbuffer == nil {
				c.MACbuffer = make(map[string][]string)
			}
			c.MACbuffer[conn.RemoteAddr().String()] = macs
			bufJSON, err := json.Marshal(c.MACbuffer)
			if err != nil {
				fmt.Println(err)
				return
			}
			s := fmt.Sprintln("JOB SPY", string(bufJSON))
			query = s
			jobInput2 <- s
			res = <-jobResult2
			fmt.Println("res:", res)
			resStr = res[conn.RemoteAddr().String()]
			resStr = resStr[len("JOB GETMAC => "):]
			fmt.Println("resStr:", resStr)
			err = json.Unmarshal([]byte(resStr), &macs)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		// do the full query since the job was accepted
		fmt.Fprintln(conn, query)

		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			panic(err)
		}
		//if len(response) >= 10 && response[:8] == "JOB SUCC" {
		respHeader, err := globals.GetHeader(response)
		if err != nil {
			panic(err)
		}
		fmt.Println("response: [", respHeader, "]")

		if respHeader == "JOB SUCC" {
			//fmt.Println("job done! result: ")
			//fmt.Println(response[len(respHeader)+1:])
			// write the response to the job result channel
			if jobResult != nil {
				data := response[len(respHeader)+1:]
				m := make(map[string]string)
				m[conn.RemoteAddr().String()] = strings.TrimSpace(fmt.Sprint(query, " => ", data))
				jobResult <- m
			}

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
