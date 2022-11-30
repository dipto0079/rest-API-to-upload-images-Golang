package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

const PORT = "8010"

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/ping", Ping).Methods("GET")
    r.HandleFunc("/upload", UploadImages).Methods("POST")

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

// handler to handle the image upload
func UploadImages(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseMultipartForm(32); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    files := r.MultipartForm.File["file"]

    var errNew string
    var http_status int

    for _, fileHeader := range files {
        // Open the file
        file, err := fileHeader.Open()
        if err != nil {
            errNew = err.Error()
            http_status = http.StatusInternalServerError
            break
        }

        defer file.Close()

        buff := make([]byte, 512)
        _, err = file.Read(buff)
        if err != nil {
            errNew = err.Error()
            http_status = http.StatusInternalServerError
            break
        }

        filetype := http.DetectContentType(buff)
        if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/jpg" {
            errNew = "The provided file format is not allowed. Please upload a JPEG,JPG or PNG image"
            http_status = http.StatusBadRequest
            break
        }

        _, err = file.Seek(0, io.SeekStart)
        if err != nil {
            errNew = err.Error()
            http_status = http.StatusInternalServerError
            break
        }

        err = os.MkdirAll("./uploads", os.ModePerm)
        if err != nil {
            errNew = err.Error()
            http_status = http.StatusInternalServerError
            break
        }

        f, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
        if err != nil {
            errNew = err.Error()
            http_status = http.StatusBadRequest
            break
        }

        defer f.Close()

        _, err = io.Copy(f, file)
        if err != nil {
            errNew = err.Error()
            http_status = http.StatusBadRequest
            break
        }
    }
    message := "file uploaded successfully"
    messageType := "S"

    if errNew != "" {
        message = errNew
        messageType = "E"
    }

    if http_status == 0 {
        http_status = http.StatusOK
    }

    resp := map[string]interface{}{
        "messageType": messageType,
        "message":     message,
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http_status)
    json.NewEncoder(w).Encode(resp)
}