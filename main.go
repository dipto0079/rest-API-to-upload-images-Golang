package main

import (
	"encoding/json"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	"test/imageUpload"
)

const PORT = "8010"

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ping", Ping).Methods("GET")
	r.HandleFunc("/upload", cusUpload).Methods("POST")

	log.Printf("Server is running on http://localhost:%s", PORT)
	log.Println(http.ListenAndServe(":"+PORT, r))
}

func Ping(w http.ResponseWriter, r *http.Request) {
	answer := map[string]interface{}{
		"messageType": "Test",
		"message":     "All OK",
		"data":        "PONG",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(answer)
}

func cusUpload(w http.ResponseWriter, r *http.Request) {

	_, js := imageUpload.UploadImages(w, r)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(js)

}
