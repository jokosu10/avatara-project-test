package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// OpenAIRequest struct for sending request to OpenAI
type OpenAIRequest struct {
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

// Function to call OpenAI API
func getOpenAIResponse(prompt string) (string, error) {
	openAIURL := "https://api.openai.com/v1/engines/davinci/completions"
	apiKey := os.Getenv("OPENAI_API_KEY")

	requestData := OpenAIRequest{
		Prompt:    prompt,
		MaxTokens: 150,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", openAIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Check if the response body is empty
	if len(body) == 0 {
		return "", fmt.Errorf("empty response from OpenAI API")
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	// Check if the expected fields are present
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("unexpected response format from OpenAI API")
	}

	text, ok := choices[0].(map[string]interface{})["text"].(string)
	if !ok {
		return "", fmt.Errorf("unable to extract text from OpenAI API response")
	}

	return text, nil
}

func chatbotHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	userPrompt := string(body)
	botResponse, err := getOpenAIResponse(userPrompt)
	if err != nil {
		http.Error(w, "Error getting response from OpenAI", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, botResponse)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error load your configuration:", err.Error())
		return
	}

	port := os.Getenv("PORT")

	fmt.Println("Starting the application...")

	http.HandleFunc("/chat", chatbotHandler)

	if err := http.ListenAndServe(":"+port, nil); err != nil { // Checked for error
		log.Fatal("The server application is error:", err.Error())
	}
}
