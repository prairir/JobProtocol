package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	creator "github.com/prairir/JobProtocol/Creator"
)

// Server struct has each channel that we need for the hanlders
// the channels are
// jobResult: string channel, read to get the result of the job
// queue: tmp channel, read the queue of connections
// job: string channel, give the job to the master loop
type Server struct {
	queueRV     chan []net.Conn
	jobInput    chan string
	jobResult   chan map[string]string
	resBuffer   map[string][]string
	queueBuffer []net.Conn
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
		var connQueue []net.Conn
		// recieves the queue
		select {
		case connQueue = <-s.queueRV:
			s.queueBuffer = connQueue
			break
		default:
			connQueue = s.queueBuffer
			log.Println("queue:", connQueue)
			break
		}
		// just init stuff
		queueJSON := make(map[string]map[string][]string)
		queueJSON["queue"] = make(map[string][]string, 0)

		// add each ip to the map
		for _, indvConn := range connQueue {
			ip := indvConn.RemoteAddr().String()
			queueJSON["queue"][ip] = make([]string, 0)
			if val, ok := s.resBuffer[ip]; ok {
				queueJSON["queue"][ip] = append(queueJSON["queue"][ip], val...)
			}
		}

		select {
		case m := <-s.jobResult:
			log.Println("got result!", m)
			for key, val := range m {
				s.resBuffer[key] = append(s.resBuffer[key], val)
				queueJSON["queue"][key] = append(queueJSON["queue"][key], val)
			}
			break
		default:
			break
		}

		w.WriteHeader(http.StatusAccepted)
		// changes the map to json
		// writes the json to the writer
		err := json.NewEncoder(w).Encode(queueJSON)

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

type jobJSON struct {
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
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}
		log.Println("form:", r.Form)
		data := make(map[string]string)
		for key := range r.Form {
			err := json.Unmarshal([]byte(key), &data)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Internal Server Error"))
				return
			}
			break
		}
		s.jobInput <- data["job"]
		w.WriteHeader(200)
		w.Write([]byte("200 - OK"))
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
		var jobResultJSON map[string]map[string]string
		jobResultJSON = make(map[string]map[string]string)

		select {
		case jobResultJSON["result"] = <-s.jobResult:
			// put the input at result value
			jobResultJSON["result"] = <-s.jobResult

			w.WriteHeader(http.StatusAccepted)
			// changes the map to json
			// writes the json to the writer
			err := json.NewEncoder(w).Encode(jobResultJSON)
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
		queueRV:   make(chan []net.Conn),
		jobResult: make(chan map[string]string, 100),
		jobInput:  make(chan string, 100),
		resBuffer: make(map[string][]string),
	}
	defer close(server.jobResult)
	defer close(server.jobInput)
	defer close(server.queueRV)

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

	creator.RunCreator(server.jobInput, server.jobResult, server.queueRV)

}
