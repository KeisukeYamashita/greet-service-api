package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/api/serviceA", func(w http.ResponseWriter, r *http.Request) {
		env := os.Getenv("ENV")

		if env == "" {
			http.Error(w, "there was no ENV defined", http.StatusInternalServerError)
			return
		}

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("bad http method"))
			return
		}

		decoder := json.NewDecoder(r.Body)

		var reqBody Message
		err := decoder.Decode(&reqBody)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("cannot decode body"))
			return
		}

		jsonBody, err := json.Marshal(reqBody)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("cannot get message "))
			return
		}

		servicebHost := os.Getenv("SERVICE_B_HOST")

		if len(servicebHost) == 0 {
			servicebHost = "localhost"
		}

		resp, err := http.Post("http://"+servicebHost+":5100/api/serviceB", "application/json", bytes.NewBuffer(jsonBody))

		if err != nil {
			http.Error(w, "Error while accessing Service B "+err.Error(), http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		msg := &Message{}

		if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err != nil {
			http.Error(w, "error while unmarshaling, err: "+err.Error(), http.StatusInternalServerError)
			return
		}

		res := &Response{
			Env:     env,
			Message: msg.Message,
		}

		response, err := json.Marshal(res)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		res := []byte("haha, i am healthy")
		w.Write(res)
	})

	http.ListenAndServe(fmt.Sprintf(":%v", "5050"), nil)
}

// Response ...
type Response struct {
	Env     string `json:"env"`
	Message string `json:"message"`
}

// Message ...
type Message struct {
	Message string `json:"message"`
}
