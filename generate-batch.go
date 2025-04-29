package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Message represents a single message in the OpenAI chat format
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Body represents the body of the chat completion request
type Body struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

// BatchRequest represents a single line in the OpenAI Batch API format
type BatchRequest struct {
	CustomID string `json:"custom_id"`
	Method   string `json:"method"`
	URL      string `json:"url"`
	Body     Body   `json:"body"`
}

// Post represents a simplified structure for CSV parsing
type Post struct {
	PostID string
}

// readPostsCSV reads the posts.csv and extracts the PostID fields
func readPostsCSV(filename string) ([]Post, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV header: %w", err)
	}

	// Find the index of post_id
	postIDIdx := -1
	for i, h := range headers {
		if h == "post_id" {
			postIDIdx = i
			break
		}
	}
	if postIDIdx == -1 {
		return nil, fmt.Errorf("post_id column not found in CSV")
	}

	var posts []Post
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		post := Post{
			PostID: record[postIDIdx],
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// readPromptTemplate reads prompt.txt into a string
func readPromptTemplate(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func main() {
	// Load prompt template
	promptTemplate, err := readPromptTemplate("prompt.txt")
	if err != nil {
		log.Fatalf("Error reading prompt.txt: %v", err)
	}

	// Load posts.csv
	posts, err := readPostsCSV("posts.csv")
	if err != nil {
		log.Fatalf("Error reading posts.csv: %v", err)
	}

	// Prepare output file (overwrite if exists)
	outputFile, err := os.Create("posts2upload.jsonl")
	if err != nil {
		log.Fatalf("Error creating posts2upload.jsonl: %v", err)
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	successCount := 0
	errorCount := 0

	for _, post := range posts {
		// Read corresponding HTML file
		htmlPath := filepath.Join("posts", post.PostID+".html")
		htmlContent, err := os.ReadFile(htmlPath)
		if err != nil {
			log.Printf("Warning: failed to read %s: %v", htmlPath, err)
			errorCount++
			continue
		}

		// Build full prompt
		fullPrompt := promptTemplate + "\n\n" + string(htmlContent)

		// Sanitize: Trim spaces
		fullPrompt = strings.TrimSpace(fullPrompt)

		// Create the batch request struct
		req := BatchRequest{
			CustomID: post.PostID,
			Method:   "POST",
			URL:      "/v1/chat/completions",
			Body: Body{
				Model: "gpt-4.1",
				Messages: []Message{
					{
						Role:    "user",
						Content: fullPrompt,
					},
				},
				MaxTokens: 2048,
			},
		}

		// Marshal to JSON
		jsonLine, err := json.Marshal(req)
		if err != nil {
			log.Printf("Warning: failed to marshal JSON for post %s: %v", post.PostID, err)
			errorCount++
			continue
		}

		// Write the JSON line
		_, err = writer.WriteString(string(jsonLine) + "\n")
		if err != nil {
			log.Printf("Warning: failed to write JSON line for post %s: %v", post.PostID, err)
			errorCount++
			continue
		}

		successCount++
	}

	// Flush any remaining buffer
	writer.Flush()

	fmt.Println("\nSummary:")
	fmt.Printf("Successfully processed %d posts.\n", successCount)
	if errorCount > 0 {
		fmt.Printf("Encountered %d errors.\n", errorCount)
	}
	fmt.Println("posts2upload.jsonl file created successfully.")
}
