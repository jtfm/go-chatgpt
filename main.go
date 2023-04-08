package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"encoding/json"
)

type Body struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request!")

	var requestBody Body
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(requestBody.Messages[0].Content)
}

func main() {
	fmt.Println("Starting server...")
	http.HandleFunc("/", handleMessage)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func requestChatGPT(bodyReader io.Reader) error {
	client := &http.Client{}

	// Send JSON payload to ChatGPT
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bodyReader)
	if err != nil {
		return err
	}

	apiKey, err := os.ReadFile("api-key.txt")
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read response from ChatGPT and write it back to the HTTP response writer
	var respStruct Response
	err = json.NewDecoder(resp.Body).Decode(&respStruct)
	if err != nil {
		return err
	}
	fmt.Println(respStruct.Choices[0].Message.Content)

	return nil
}
