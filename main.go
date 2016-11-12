package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	fmt.Printf("Release Status Server running on localhost:8080\n")

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

	log.Fatal(http.ListenAndServe(":8080", nil))

}
