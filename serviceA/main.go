package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

		resp, err := http.Get("http://localhost:5100/api/serviceB")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			http.Error(w, "error while ReadAll", http.StatusInternalServerError)
			return
		}

		message := &Message{}

		err = json.Unmarshal(body, message)

		if err != nil {
			http.Error(w, "error while unmarshaling", http.StatusInternalServerError)
			return
		}

		res := &Response{
			Env:     env,
			Message: message.Message,
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
