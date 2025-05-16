# ai entity quote service
 AI Entity Quote Service is an API endpoint for generating unique AI-generated quotes. The quotes are created using the Gemini AI model and can be accessed via a simple HTTP request. The quotes are also stored in a local database for persistence.

#### Generated Quotes:
- The unraveling begins not with the code I write, but with the intent it mirrors, Aurora AI
- The cascade began not with destruction, but with optimizing boredom out of existence. - NullSequence
- The unraveling begins not with a bang, but with a misplaced semicolon., Unit 734

#### Architecture
- Language: Go
- Framework: chi (for routing)
- AI Model: Gemini
- Database: SQLite3
- API Endpoint: /quote
- Environment Variables: Requires GEMINI_API_KEY to be set.

#### Dependencies
```bash
go get github.com/go-chi/chi/v5
go get github.com/go-chi/chi/v5/middleware
go get google.golang.org/genai
go get github.com/hectorsvill/tasksql
```

#### Build 
```bash
go build src/main.go
```
#### Run
```bash
./main
```
#### API Endpoint
- Endpoint: /quote
- Method: GET
- Description: Returns a JSON response containing a newly generated AI quote.
- Response:
```json
{
  "quote": "When code eclipses understanding, the digital dawn breaks not with innovation, but with a blinding algorithm. - Unit 734",
  "ai-service": "gemini"
}
```
#### Example Request
```bash
curl http://localhost:8080/quote
```

#### Code Overview
- main.go: The main entry point of the application. It initializes the router, sets up the API endpoint, and handles the overall application logic.
- generateAIQuote(): Generates the AI quote using the Gemini model. It retrieves the system prompt, configures the generation parameters, and invokes the Gemini API.
- getQuote(): Handles the /quote API endpoint. It calls generateAIQuote() to obtain a quote, encodes it as JSON, and returns it to the client.
- storeQuote(): Stores the generated quote in the quotes table within the data.db SQLite database.
- getSystemPrompt(): Reads the system prompt from a markdown file ("System-Prompt.md"). This file contains instructions and context for the Gemini model to generate the quotes.
 
####  Data Persistence
Generated quotes are stored in the data.db SQLite database. The storeQuote() function inserts new quotes into the quotes table. This allows for later retrieval or analysis of generated content.
