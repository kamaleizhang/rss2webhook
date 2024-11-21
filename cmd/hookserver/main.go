package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlePostRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var jsonData interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	jsonOutput, err := json.MarshalIndent(jsonData, "", "    ")
	if err != nil {
		http.Error(w, "Error formatting JSON for output", http.StatusInternalServerError)
		return
	}
	log.Printf("%v", string(jsonOutput))

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonOutput)
	if err != nil {
		return
	}
}
