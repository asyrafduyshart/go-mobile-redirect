package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// RedirectType data to redirect
type RedirectType struct {
	Mobile  string `json:"mobile"`
	Desktop string `json:"desktop"`
}

// URLRedirect data to redirect
var URLRedirect map[string]RedirectType

func isMobile(useragent string) bool {
	// the list below is taken from
	// https://github.com/bcit-ci/CodeIgniter/blob/develop/system/libraries/User_agent.php

	mobiles := []string{"Mobile Explorer", "Palm", "Motorola", "Nokia", "Palm", "Apple iPhone", "iPad", "Apple iPod Touch", "Sony Ericsson", "Sony Ericsson", "BlackBerry", "O2 Cocoon", "Treo", "LG", "Amoi", "XDA", "MDA", "Vario", "HTC", "Samsung", "Sharp", "Siemens", "Alcatel", "BenQ", "HP iPaq", "Motorola", "PlayStation Portable", "PlayStation 3", "PlayStation Vita", "Danger Hiptop", "NEC", "Panasonic", "Philips", "Sagem", "Sanyo", "SPV", "ZTE", "Sendo", "Nintendo DSi", "Nintendo DS", "Nintendo 3DS", "Nintendo Wii", "Open Web", "OpenWeb", "Android", "Symbian", "SymbianOS", "Palm", "Symbian S60", "Windows CE", "Obigo", "Netfront Browser", "Openwave Browser", "Mobile Explorer", "Opera Mini", "Opera Mobile", "Firefox Mobile", "Digital Paths", "AvantGo", "Xiino", "Novarra Transcoder", "Vodafone", "NTT DoCoMo", "O2", "mobile", "wireless", "j2me", "midp", "cldc", "up.link", "up.browser", "smartphone", "cellphone", "Generic Mobile"}

	for _, device := range mobiles {
		if strings.Index(useragent, device) > -1 {
			return true
		}
	}
	return false
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
	host := r.Host
	if isMobile(ua) {
		http.Redirect(w, r, "https://"+URLRedirect[host].Mobile, 302)
		return
	}
	http.Redirect(w, r, "https://"+URLRedirect[host].Desktop, 302)
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

	http.HandleFunc("/", redirect)
	http.HandleFunc("/set-to-newest-data", getNewestData)
	err := http.ListenAndServe(":"+getEnv("PORT", "9090"), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
