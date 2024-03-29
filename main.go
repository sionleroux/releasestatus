// The releasestatus binary is a small daemon to track the state of QA releases.
package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// TODO: make this non-global
var Lock sync.Mutex // synchronise access to the release status

func main() {
	// server name is omitted before : because it's localhost
	host := fmt.Sprint(":", httpPort())
	log.Printf("Release Status Server running on localhost%v\n", host)

	// current release (starts off empty)
	cur := &Release{}

	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, buildResponse(cur.start(r.URL.Query().Get("name"))))
	})
	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, buildResponse(cur.stop()))
	})
	http.HandleFunc("/block", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, buildResponse(cur.block()))
	})
	http.HandleFunc("/unblock", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, buildResponse(cur.unblock()))
	})
	log.Fatal(http.ListenAndServe(host, nil))
}

// buildResponse formats the result of an action in a way it can be consumed by
// an API client.
func buildResponse(err error) int {
	if err != nil {
		return 0
	}
	return 1
}

// Release persists meta-data about a running release.
type Release struct {
	Author    string    // the person who started the release
	Running   bool      // whether the release is running or not
	Blocked   bool      // whether further releases are allowed or not
	StartedAt time.Time // timestamp when the release started
}

// start marks a Release as started unless another release is already running.
func (r *Release) start(author string) error {
	if author == "" {
		author = "an unknown user"
	}
	Lock.Lock()
	defer Lock.Unlock()
	if r.Running {
		log.Printf("Refusing start request because release already started by %v at %v\n", r.Author, r.StartedAt)
		return errors.New("release already in progress")
	}
	if r.Blocked {
		log.Printf("Refusing start request because releasing is blocked")
		return errors.New("release blocked")
	}
	log.Printf("Starting new release by %v\n", author)
	*r = Release{author, true, false, time.Now()}
	return nil
}

// stop marks a Release as stopped unless it is already stopped.
func (r *Release) stop() error {
	Lock.Lock()
	defer Lock.Unlock()
	if !r.Running {
		log.Print("Refusing to stop release because no release running")
		return errors.New("no release to stop")
	}
	log.Print("Stopping release")
	blocked := r.Blocked
	*r = Release{}
	r.Blocked = blocked // reset blocked status
	return nil
}

// block marks a Release as blocked unless it is already blocked.
func (r *Release) block() error {
	Lock.Lock()
	defer Lock.Unlock()
	if r.Blocked {
		log.Print("Refusing to block release because release already blocked")
		return errors.New("release already blocked")
	}
	log.Print("Blocking release")
	r.Blocked = true
	return nil
}

// unblock marks a Release as unblocked unless it is already unblocked.
func (r *Release) unblock() error {
	Lock.Lock()
	defer Lock.Unlock()
	if !r.Blocked {
		log.Print("Refusing to unblock release because release is not blocked")
		return errors.New("release already unblocked")
	}
	log.Print("Unblocking release")
	r.Blocked = false
	return nil
}

// Gets program port from $RS_PORT env var.
func httpPort() int {
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
