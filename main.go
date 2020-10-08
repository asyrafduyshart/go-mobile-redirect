package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/mssola/user_agent"
)

// RedirectType data to redirect
type RedirectType struct {
	Mobile  string `json:"mobile"`
	Desktop string `json:"desktop"`
}

// URLRedirect data to redirect
var URLRedirect map[string]RedirectType

func isMobile(useragent string) bool {
	return user_agent.New(useragent).Mobile()
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func redirect(w http.ResponseWriter, r *http.Request) {
	ua := r.Header.Get("User-Agent")

	uri := r.RequestURI

	host := r.Host
	if isMobile(ua) {
		http.Redirect(w, r, "https://"+URLRedirect[host].Mobile+uri, 302)
		return
	}
	http.Redirect(w, r, "https://"+URLRedirect[host].Desktop+uri, 302)
}

func getNewestData(w http.ResponseWriter, r *http.Request) {
	err := setNewHost()
	if err != nil {
		fmt.Println("Failed to set new data host")
		w.WriteHeader(500)
		w.Write([]byte("Failed to setd ata"))
		return
	}
	fmt.Println("Data successfully fetched")
	w.WriteHeader(200)
	w.Write([]byte("Data successfully fetched"))
	return
}

func setNewHost() error {
	resp, err := http.Get(getEnv("REDIRECT_DATA_URL", "https://some-url-containing-example.redirect.json"))
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&URLRedirect)
	if err != nil {
		initialize()
	}
	return err
}

func initialize() {
	jsonFile, err := os.Open("redirect.json")
	if err != nil {
		log.Fatal("Error reading redirect.json: ", err)
	}
	fmt.Println("Successfully Opened redirect.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &URLRedirect)
}

func main() {
	setNewHost()
	// initialize()

	http.HandleFunc("/", redirect)
	http.HandleFunc("/set-to-newest-data", getNewestData)
	err := http.ListenAndServe(":"+getEnv("PORT", "9090"), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
