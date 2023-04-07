package main

import (
	"os"
    "fmt"
    "net/http"
    "io/ioutil"
	"bytes"
	"encoding/json"
)

type Body struct {
	Model string `json:"model"`
	Messages []Messages `json:"messages"`
}

type Messages struct {
	Role string `json:"role"`
	Content string `json:"content"`
}

func handleMessage(w http.ResponseWriter, r *http.Request) {

	body := Body{
		Model: "gpt-3.5-turbo",
		Messages: []Messages{
			{
				Role: "user",
				Content: "Hello",
			},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bodyReader := bytes.NewReader(jsonBody)

	client := &http.Client{}

    // Send JSON payload to ChatGPT
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bodyReader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiKey, err := os.ReadFile("api-key.txt")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := client.Do(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // Read response from ChatGPT and write it back to the HTTP response writer
    response, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, string(response))
}

func main() {
    http.HandleFunc("/", handleMessage)
    http.ListenAndServe(":8080", nil)
}