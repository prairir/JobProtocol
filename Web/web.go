package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	creator "github.com/prairir/JobProtocol/Creator"
)

// struct has each channel that we need for the hanlders
// the channels are
// jobResult: string channel, read to get the result of the job
// queue: tmp channel, read the queue of connections
// job: string channel, give the job to the master loop
type Server struct {
	jobResult chan string
	queueRV   chan []net.Conn
	queueTR   chan int
	jobInput  chan string
	connQueue []net.Conn
}

// go handler for queue GET request
// method: GET
// output:
// ```
// {
// queue: [
// 		"123.321.123.0",
// 		"123.321.323.0",
// 		"123.321.333.0",
//		]
//	}
// ```
func (s *Server) queueHandler(w http.ResponseWriter, r *http.Request) {
	// if method is GET
	if r.Method == http.MethodGet {

		// send the message to write to it
		s.queueTR <- 1

		var connQueue []net.Conn
		// recieves the queue
		select {
		case connQueue = <-s.queueRV:
			break
		default:
			w.Write([]byte("{\"queue\": [ ]}"))
			return
		}

		if connQueue == nil {
			connQueue = s.connQueue
		}

		// just init stuff
		var queueJson map[string][]string
		queueJson = make(map[string][]string)
		queueJson["queue"] = make([]string, 0)

		// add each ip to the map
		for _, indvConn := range connQueue {
			queueJson["queue"] = append(queueJson["queue"], indvConn.RemoteAddr().String())
		}

		w.WriteHeader(http.StatusAccepted)
		// changes the map to json
		// writes the json to the writer
		err := json.NewEncoder(w).Encode(queueJson)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Bad request"))
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Wrong request method"))
	}

}

type jobJson struct {
	job string
}

// go handler for job POST request
// method: POST
// data input:
// ```
// {
// 	"job": "JOB EQN 1+2"
// }
// ```
// output: status code
func (s *Server) jobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)

		var data jobJson
		err := decoder.Decode(&data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Bad request"))
			return
		}

		s.jobInput <- data.job
		w.WriteHeader(200)
		w.Write([]byte("200 - Ok response"))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Wrong request method"))
	}
}

// go handler for job GET request
// method: GET
// output:
// ```
// {
// 	result: "2/32"
// }
// ```
func (s *Server) jobResultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// create the map
		var jobResultJson map[string]string
		jobResultJson = make(map[string]string)

		select {
		case jobResultJson["result"] = <-s.jobResult:
			// put the input at result value
			jobResultJson["result"] = <-s.jobResult

			w.WriteHeader(http.StatusAccepted)
			// changes the map to json
			// writes the json to the writer
			err := json.NewEncoder(w).Encode(jobResultJson)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("400 - Bad request"))
				return
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Bad request"))
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Wrong request method"))
	}
}

func main() {
	//file server to return our files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// creating a server struct so we can share channels to and from the
	server := &Server{
		jobResult: make(chan string),
		queueRV:   make(chan []net.Conn),
		queueTR:   make(chan int, 100000),
		jobInput:  make(chan string, 100000),
	}

	defer close(server.jobResult)
	defer close(server.jobInput)
	defer close(server.queueRV)
	defer close(server.queueTR)
	//handlers
	http.HandleFunc("/api/queue", server.queueHandler)
	http.HandleFunc("/api/job", server.jobHandler)
	http.HandleFunc("/api/jobResult", server.jobResultHandler)

	//run the stuff
	log.Println("Listening on http://localhost:8080...")
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	creator.RunCreator(server.queueTR, server.queueRV, server.jobInput, server.jobResult)

}
