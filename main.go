package main

// Magic string to reload using reflex:
// reflex -g '*.go' -s -- sh -c 'go build && ./netlify-teams-webhook'

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
)

// dump request to stdout
func dump(w http.ResponseWriter, req *http.Request) {

	// Save a copy of this request for debugging.
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Request dump:")
	fmt.Println("-----------")
	fmt.Println(string(requestDump))
	fmt.Println("-----------")

}

// getPort returns port from environment variable PORT if set, otherwise return
// defaultPort
func getPort(defaultPort int) int {
	if value, ok := os.LookupEnv("PORT"); ok {
		port, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		return port
	}
	return defaultPort
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/dump", dump)

	port := ":" + strconv.Itoa(getPort(8090))
	log.Println(fmt.Sprintf("Server running on http://localhost%s 🐹", port))
	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("could not run the server %v", err)
		return
	}
}
