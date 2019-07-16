package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()

	http.HandleFunc("/api/serviceA", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("handling request")
		env := os.Getenv("ENV")

		if env == "" {
			http.Error(w, "there was no ENV defined", http.StatusInternalServerError)
			return
		}

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			logger.Error(fmt.Sprintf("bad http request method:%s", r.Method))
			w.Write([]byte("bad http method"))
			return
		}

		decoder := json.NewDecoder(r.Body)

		var reqBody Message
		err := decoder.Decode(&reqBody)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("bad http request body error:%v", err))
			w.Write([]byte("cannot decode body"))
			return
		}

		jsonBody, err := json.Marshal(reqBody)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("bad response from service B error:%v", err))
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

		logger.Info("got message form service B")
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
