package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/genai"
	"github.com/google/uuid"

	"github.com/hectorsvill/tasksql"
)

func textToImage(prompt string) string {
	apiKey := getAPIKey()
	ctx := context.Background()

	client, _ := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
		Backend: genai.BackendGeminiAPI,
	})

	config := &genai.GenerateContentConfig{
		ResponseModalities: []string{"TEXT", "IMAGE"},
	}
	
	result, _ := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash-preview-image-generation",
		genai.Text(prompt),
		config,
	)
	id := uuid.New().String()
	for _, part := range result.Candidates[0].Content.Parts {
		if part.Text != "" {
			fmt.Println(part.Text)
		} else if part.InlineData != nil {
			imageBytes := part.InlineData.Data
			outputFilename := fmt.Sprintf("./src/img/%s.png", id)
			_ = os.WriteFile(outputFilename, imageBytes, 0644)
		}
  	}
	return id
}

func getAPIKey() string {
	apikey := os.Getenv("GEMINI_API_KEY")
	if apikey == "" {
		log.Fatal("Please set the GEMINI_API_KEY")
	}
	return apikey
}

func storeQuote(q string) {
	taskSql, err := tasksql.NewDB("data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer taskSql.Close()

	taskSql.CreateTableIfNotExist("quotes")
	taskSql.Post("quotes", q)
}

func getSystemPrompt() string {
	filePath := "./src/System-Prompt.md"
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	systemPrompt := string(content)
	return systemPrompt
}

func generateAIQuote() string {
	ctx := context.Background()
	apikey := getAPIKey()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apikey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}
	prompt := getSystemPrompt()
	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(prompt, genai.RoleUser),
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text("Create a unique ai quote about dooms day,hacking, building robots,building ai's,coding or anything else interesting in tech and tag it by the ai author. Split the quote and auther with a comma"),
		config,
	)
	if err != nil {
		log.Fatal(err)
	}
	go storeQuote(result.Text())
	log.Println(result.Text())
	return result.Text()
}

func getQuote(w http.ResponseWriter, r *http.Request) {
	quote := generateAIQuote()
	id := textToImage(quote)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"quote":      quote,
		"ai-service": "gemini",
		"id": id,
	})
	
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/quote", getQuote)

	fmt.Printf("Starting Server on :8080\n")
	fmt.Printf("Endpoint: :8080/quote")

	http.ListenAndServe(":8080", router)
}
