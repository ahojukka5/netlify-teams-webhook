package main

// Magic string to reload using reflex:
// reflex -g '*.go' -s -- sh -c 'go build && ./netlify-teams-webhook'

import (
	"encoding/json"
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

// LogAccessAttributes is a part of NetlifyPayload
type LogAccessAttributes struct {
	Type     string `json:"type"`
	URL      string `json:"url"`
	Endpoint string `json:"endpoint"`
	Path     string `json:"path"`
	Token    string `json:"token"`
}

// NetlifyPayload is struct for json data what Netlify sends
type NetlifyPayload struct {
	ID                  string              `json:"id"`
	SiteID              string              `json:"site_id"`
	BuildID             string              `json:"build_id"`
	DeployURL           string              `json:"deploy_url"`
	DeploySSLURL        string              `json:"deploy_ssl_url"`
	CreatedAt           string              `json:"created_at"`
	UpdatedAt           string              `json:"updated_at"`
	PublishedAt         string              `json:"published_at"`
	UserID              string              `json:"user_id"`
	CommitRef           string              `json:"commit_ref"`
	Branch              string              `json:"branch"`
	LogAccessAttributes LogAccessAttributes `json:"log_access_attributes"`
}

// Fact is part of Card
type Fact struct {
	Name  string `json:"name"`
	Value string `json:"values"`
}

// Section is part of Card
type Section struct {
	ActivityTitle string `json:"activityTitle"`
	Facts         []Fact `json:"facts"`
}

// Card contains the payload sent to MS Teams
type Card struct {
	Type       string    `json:"@type"`
	Context    string    `json:"@context"`
	ThemeColor string    `json:"themeColor"`
	Summary    string    `json:"summary"`
	Sections   []Section `json:"sections"`
}

func deployCreated(w http.ResponseWriter, req *http.Request) {

	// Check header
	if req.Header.Get("X-Netlify-Event") != "deploy_created" {
		msg := "X-Netlify-Event header is not deploy_created"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	// Declare a new NewlifyPayload struct.
	var payload NetlifyPayload

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Do something with the Person struct...
	fmt.Fprintf(w, "NetlifyPayload: %+v", payload)
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
	mux.HandleFunc("/deploy_created", deployCreated)

	port := ":" + strconv.Itoa(getPort(8090))
	log.Println(fmt.Sprintf("Server running on http://localhost%s 🐹", port))
	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("could not run the server %v", err)
		return
	}
}
