package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {

	// server name is omitted before : because it's localhost
	host := fmt.Sprint(":", getPort())
	log.Printf("Release Status Server running on localhost%v\n", host)

	releasing := false

	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		if releasing {
			log.Print("Refusing start request because release already running")
			fmt.Fprint(w, "0")
		} else {
			log.Print("Starting new release")
			releasing = true
			fmt.Fprint(w, "1")
		}
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		if releasing {
			log.Print("Stopping release")
			releasing = false
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
	fmt.Println(port)

	return port
}
