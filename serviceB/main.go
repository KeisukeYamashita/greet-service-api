package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/api/serviceB", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("bad http method"))
			return
		}

		decoder := json.NewDecoder(r.Body)

		var msg Message
		err := decoder.Decode(&msg)

		var resMsg string

		switch msg.Message {
		case "Hi":
			resMsg = "Hi"
		case "Bye":
			resMsg = "Bye"
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		res := &Message{
			Message: os.Getenv("SECRET_MESSAGE_PREFIX") + " " + resMsg,
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

	http.ListenAndServe(fmt.Sprintf(":%v", "5100"), nil)
}

// Message ...
type Message struct {
	Message string `json:"message"`
}
