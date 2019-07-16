package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logger, _ := zap.NewProduction()

	http.HandleFunc("/api/serviceB", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("handling request")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			logger.Error(fmt.Sprintf("bat http request method:%s", r.Method))
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
			logger.Error(fmt.Sprintf("failed to marshal responce error:%v", err))
			return
		}

		logger.Info("responding to service A")
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
