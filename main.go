package main

// Magic string to reload using reflex:
// reflex -g '*.go' -s -- sh -c 'go build && ./netlify-teams-webhook'

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
)

const teamsWebhookURLEnv string = "TEAMS_WEBHOOK_URL"

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

// Links contains permalink
type Links struct {
	Permalink string `json:"permalink"`
	Alias     string `json:"alias"`
}

// NetlifyPayload is struct for json data what Netlify sends
type NetlifyPayload struct {
	ID                  string              `json:"id"`
	SiteID              string              `json:"site_id"`
	BuildID             string              `json:"build_id"`
	Name                string              `json:"name"`
	DeployURL           string              `json:"deploy_url"`
	DeploySSLURL        string              `json:"deploy_ssl_url"`
	CreatedAt           string              `json:"created_at"`
	UpdatedAt           string              `json:"updated_at"`
	PublishedAt         string              `json:"published_at"`
	DeployTime          int                 `json:"deploy_time"`
	UserID              string              `json:"user_id"`
	CommitRef           string              `json:"commit_ref"`
	Branch              string              `json:"branch"`
	LogAccessAttributes LogAccessAttributes `json:"log_access_attributes"`
	Links               Links               `json:"links"`
}

// Fact is part of Card
type Fact struct {
	Name  string `json:"name"`
	Value string `json:"value"`
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
	// fmt.Fprintf(w, "NetlifyPayload: %+v", payload)

	title := payload.Name + " published new page at " + payload.PublishedAt

	permalink := Fact{
		Name:  "permalink",
		Value: fmt.Sprintf("<a href=\"%s\">%s<a>", payload.Links.Permalink, payload.Links.Permalink),
	}

	deployURL := Fact{
		Name:  "deploy_url",
		Value: fmt.Sprintf("<a href=\"%s\">%s<a>", payload.DeploySSLURL, payload.DeploySSLURL),
	}

	buildID := Fact{
		Name:  "build_id",
		Value: payload.BuildID,
	}

	createdAt := Fact{
		Name:  "created_at",
		Value: payload.CreatedAt,
	}

	publishedAt := Fact{
		Name:  "published_at",
		Value: payload.PublishedAt,
	}

	deployTime := Fact{
		Name:  "deploy_time",
		Value: strconv.Itoa(payload.DeployTime),
	}

	facts := []Fact{permalink, deployURL, buildID, createdAt, publishedAt, deployTime}

	cardSection := Section{
		ActivityTitle: title,
		Facts:         facts,
	}

	cardSections := []Section{cardSection}

	card := Card{
		Type:       "MessageCard",
		Context:    "http://schema.org/extensions",
		ThemeColor: "0076D7",
		Summary:    title,
		Sections:   cardSections,
	}
	cardJSON, err := json.Marshal(&card)
	if err != nil {
		fmt.Println(err)
		return
	}

	teamsWebhookURL := os.Getenv(teamsWebhookURLEnv)

	fmt.Println("Sending " + string(cardJSON) + "\n\nto " + teamsWebhookURL + "\n")
	request, err := http.NewRequest("POST", teamsWebhookURL, bytes.NewBuffer(cardJSON))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	fmt.Println("Response Status:", response.Status)
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
	teamsWebhookURL, exists := os.LookupEnv(teamsWebhookURLEnv)
	if !exists {
		println("Environment variable " + teamsWebhookURLEnv + " must be set.")
		return
	}
	fmt.Println("Sending Netlify notification to MS Teams Webhook at " + teamsWebhookURL)

	mux := http.NewServeMux()
	mux.HandleFunc("/dump", dump)
	mux.HandleFunc("/deploy_created", deployCreated)

	port := ":" + strconv.Itoa(getPort(8090))
	fmt.Println(fmt.Sprintf("Server running on http://localhost%s 🐹", port))
	err := http.ListenAndServe(port, mux)
	if err != nil {
		fmt.Println("could not run the server", err)
		return
	}
}
