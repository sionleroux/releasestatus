package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	// server name is omitted before : because it's localhost
	host := fmt.Sprint(":", getPort())
	log.Printf("Release Status Server running on localhost%v\n", host)

	// current release (starts off empty)
	cur := &Release{}

	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, buildResponse(cur.start(r.URL.Query().Get("name"))))
	})
	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, buildResponse(cur.stop()))
	})
	log.Fatal(http.ListenAndServe(host, nil))
}

func buildResponse(err error) int {
	if err != nil {
		return 0
	}
	return 1
}

// Persistent release meta-data
type Release struct {
	Author    string    // the person who started the release
	Running   bool      // whether the release is running or not
	StartedAt time.Time // timestamp when the release started
}

// start marks a release as started unless another release is already running
func (r *Release) start(author string) error {
	if author == "" {
		author = "an unknown user"
	}
	if r.Running {
		log.Printf("Refusing start request because release already started by %v at %v\n", r.Author, r.StartedAt)
		return errors.New("release already in progress")
	}
	log.Printf("Starting new release by %v\n", author)
	*r = Release{author, true, time.Now()}
	return nil
}

func (r *Release) stop() error {
	if !r.Running {
		log.Print("Refusing to stop release because no release running")
		return errors.New("no release to stop")
	}
	log.Print("Stopping release")
	*r = Release{}
	return nil
}

// Gets program port from $RS_PORT env var
func getPort() int {
	const defaultport string = "8080"
	portstr := os.Getenv("RS_PORT")
	if portstr == "" {
		log.Printf("No port specified, falling back to default: %s\n", defaultport)
		portstr = defaultport
	}
	port, err := strconv.Atoi(portstr)
	if err != nil {
		log.Fatal("Port is non-numeric: ", err)
	}
	return port
}
