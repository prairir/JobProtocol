package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

type IdvSession struct {
	id         int64
	state      int
	connection net.Conn
}

func main() {
	// read IP from console
	reader1 := bufio.NewReader(os.Stdin)
	fmt.Print("Enter IP:")
	IP, _ := reader1.ReadString('\n')

	// read port from consol
	reader2 := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Port:")
	port, _ := reader2.ReadString('\n')

	// add the IP and port, remove \n
	addr := IP[:len(IP)-1] + ":" + port[:len(port)-1]
	fmt.Print(addr)

	// open port to tcp connection
	tcpAddr, err := net.ResolveTCPAddr("tcp", tcpAddr)
	checkError(err)

	// create listener on that open port
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// the mutex makes it memory safe
	var mutex sync.mutex
	queue := list.new()

	for {
		conn, err := listener.Accept()
		if err != nil {
			closeConnection(conn)
		}
	}

	go handleConnection(conn, &mutex, queue)
}

func handleConnections(conn net.Conn, mutex *sync.mutex, queue *list.list) {
	// state values
	// 0 waiting for connection
	// 1 connection established
	// 2 job available
	// 3 send job
	// 4 closed

	state := 0

	currSession := IdvSession{
		// session ID is time in nano seconds, similar to seeker
		time.Now().UnixNano(),
		state,
		conn,
	}

	// adds new element, currSession, with unique struct values to back of queue
	mutex.Lock()
	queue.Pushback(currSession)
	mutex.Unlock()

	//Event handling with for loop
	for {
		// checks if connection contains an error
		result, err := ioutil.ReadAll(conn)
		if err != nil {
			closeConnection(conn)
		}

		// cleans result (byte slice) into string by trimming white space
		cleanedResult := strings.TrimSpace(string(result))

		var positiong int = getPosition(&currSession, queue)

		// if state is 0, look for connection
		if state == 0 && strings.Compare(cleanedResult, "HELLOACK") == 0 {
			state = 1
		} else if state == 0 {
		    conn.Write([]byte("HELLO"))
			state = 0
		}

		//hello is acknowledged, connection established
		if state == 1 {
			// check position of seeker and avl message
			if position == 0 {
				conn.Write([]byte("AVL"))

				// checks if seeker responds with avlack or full
				if strings.Compare(cleanedResult, "AVLACK") == 0 {
					// if AVLACK, advance to state 2
					state = 2
				} else {
					// if seeker is full, go back to state 1
					// try to re establish availability
					state = 1
				}
			}
		}

		// seeker is available
		if state == 2 {
			conn.Write([]byte("JOB TIME"))

			// if seeker responds with done time
			if strings.Compare(cleanedResult[:9], "DONE TIME") == 0 {
				state = 3
			}
		}

		// send seeker an equation (job)
		if state == 3 {
			var eq String = "1 + 2"
			conn.Write([]byte("JOB EQ" + eq))

			if strings.Compare(cleanedResult[:6], "DONE EQ") == 0 {
				state = 4
			}
		}

		// end at state 4
		if state == 4 {
			break
		}
	}
}

func getPosition(currSession *IdvSession, mutex *sync.mutex, queue *list.list) int {
	mutex.Lock()
	// interate throughh list to print contents
	position := 0

	for e, index := queue.Front(), 0; e != nil; index = e.Next(), index + 1 {
		fmt.Println(e.Value)

		//compare the value of currSession to session in queue
		if currSession.id == e.Value.(*IdvSession).id {
			position = index
			break
		}
	}
	mutex.Unlock()
	return position
}

func closeConnection(conn net.Conn) {
	conn.Write([]Byte("Closing connection with job creator."))
	conn.Close()
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
