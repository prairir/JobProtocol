package main

import (
	"encoding/json"
	"fmt"
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
	job       chan string
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

		// recieves the queue
		connQueue := <-s.queueRV

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
		json.NewEncoder(w).Encode(queueJson)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Wrong request method"))
	}

}

func (s *Server) jobHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("URL is: ", r.URL.Path)
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func (s *Server) jobResultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("URL is: ", r.URL.Path)
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path)
}

func main() {
	//file server to return our files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// creating a server struct so we can share channels to and from the
	server := &Server{
		jobResult: make(chan string),
		queueRV:   make(chan []net.Conn),
		queueTR:   make(chan int),
		job:       make(chan string),
	}

	//handlers
	http.HandleFunc("/api/queue", server.queueHandler)
	http.HandleFunc("/api/job", server.jobHandler)
	http.HandleFunc("/api/jobResult", server.jobResultHandler)

	//run the stuff
	log.Println("Listening on http://localhost:8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

	creator.RunCreator(server.queueTR, server.queueRV)

}
