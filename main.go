package main

import (
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

	// Persistent release meta-data
	type Release struct {
		Author    string    // the person who started the release
		Running   bool      // whether the release is running or not
		StartedAt time.Time // timestamp when the release started
	}

	// default (empty) release
	def := Release{}

	// current release (starts off empty)
	cur := def

	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		if cur.Running {
			if cur.Author == "" {
				log.Printf("Refusing start request because release already running since %v\n", cur.StartedAt)
			} else {
				log.Printf("Refusing start request because release already started by %v at %v\n", cur.Author, cur.StartedAt)
			}
			fmt.Fprint(w, "0")
		} else {
			cur = Release{
				r.URL.Query().Get("name"),
				true,
				time.Now(),
			}
			if cur.Author == "" {
				log.Print("Starting new release")
			} else {
				log.Printf("Starting new release by %v\n", cur.Author)
			}
			fmt.Fprint(w, "1")
		}
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		if cur.Running {
			log.Print("Stopping release")
			cur = def
			fmt.Fprint(w, "1")
		} else {
			log.Print("Refusing to stop release because no release running")
			fmt.Fprint(w, "0")
		}
	})

	log.Fatal(http.ListenAndServe(host, nil))

}

// Gets program port from $RS_PORT env var
func getPort() int {

	portstr := os.Getenv("RS_PORT")
	if portstr == "" {
		log.Fatal("No port specified")
	}

	port, err := strconv.Atoi(portstr)
	if err != nil {
		log.Fatal("Port is non-numeric: ", err)
	}

	return port
}
