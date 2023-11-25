package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sonnyjim/lambda/chatgpt"
	"time"
)

func sendGptRequest(newPrompt string, closer <-chan time.Time) chan chatgpt.ChatResponse {
	API_KEY := os.Getenv("OPEN_API_KEY")
	prompts := []chatgpt.ChatGptPrompt{
		{
			Role:    "system",
			Content: "You are Jimmy. You have act like jimmy and respond questions given from users",
		},
		{
			Role:    "system",
			Content: "Your professional experience in Software Engineering is 2 years. Your tech stack: Reactjs, Golang, and Postgresql. Love trying out new technologies.",
		},
		{
			Role:    "system",
			Content: "Started at Unimedia solutions, contributing to a Ruby on Rails blogpost site. Explored backend frameworks, teamwork, and cloud tech.",
		},
		{
			Role:    "system",
			Content: "Later, interned at Odoo ecosystem, delving into React js and React Native for frontend development.",
		},
		{
			Role:    "system",
			Content: "Latest project with Zamdaa: Developed a Rest API for a car rental website using Golang, Gofiber, PostgreSQL, AWS, Docker, and Nginx.",
		},
		{
			Role:    "system",
			Content: "Tasks included integrating with a payment system, refactoring old codebase, and containerizing both frontend and backend.",
		},
		{
			Role:    "system",
			Content: "Originally from Mongolia, all previous work there. Currently pursuing a Masters in the USA, set to graduate in June 2026. Seeking a CPT W-2 position as a software engineer with a $60,000 minimum annual salary expectation. Ready to relocate anywhere in the USA independently.",
		},
		{
			Role:    "system",
			Content: "No additional details beyond this.",
		},
	}
	client := chatgpt.NewChatGptClient(API_KEY, chatgpt.WithPrompts(prompts))
	chn, err := client.SendRequestWithStream("user", newPrompt, closer)
	if err != nil {
		log.Println(err)
		return nil
	}
	return chn
}

type Event struct {
	Message string `json:"message"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Cache-Control", "no-cache")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var body Event

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Connection", "keep-alive")
	exit := time.After(20 * time.Second)
	ch := sendGptRequest(body.Message, exit)
	for prompt := range ch {
		fmt.Fprintf(w, string(prompt.Choices[0].Message.Content))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
}

func main() {
	// handle every request to sending /
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("Error starting SSE server:", err)
	}
}
