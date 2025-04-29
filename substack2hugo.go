package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Post represents a single blog post entry from the CSV
type Post struct {
	PostID      string
	PostDate    string
	IsPublished string
	Title       string
}

// SEOContent represents the parsed content from batch_output.jsonl
type SEOContent struct {
	Title       string
	Description string
	Keywords    string
}

// OutputRecord represents a line in batch_output.jsonl. It's a nested struct, thus the other ones.
type OutputRecord struct {
	CustomID string     `json:"custom_id"`
	Response Response   `json:"response"`
}

type Response struct {
	Body       Body   `json:"body"`
}

type Body struct {
	Choices          []Choice  `json:"choices"`
}

type Choice struct {
	Index        int      `json:"index"`
	Message      Message  `json:"message"`
}

type Message struct {
	Content     string        `json:"content"`
}

// readPostsCSV reads and parses the posts.csv file
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

	fieldIdx := make(map[string]int)
	for i, h := range headers {
		fieldIdx[h] = i
	}

	required := []string{"post_id", "post_date", "is_published", "title"}
	for _, key := range required {
		if _, ok := fieldIdx[key]; !ok {
			return nil, fmt.Errorf("missing required field: %s", key)
		}
	}

	var posts []Post
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		post := Post{
			PostID:      record[fieldIdx["post_id"]],
			PostDate:    record[fieldIdx["post_date"]],
			IsPublished: record[fieldIdx["is_published"]],
			Title:       record[fieldIdx["title"]],
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// cleanOrCreateFolder deletes all files inside folder or creates it if missing
func cleanOrCreateFolder(folder string) error {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		return os.Mkdir(folder, 0755)
	}

	entries, err := os.ReadDir(folder)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		err = os.RemoveAll(filepath.Join(folder, entry.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

// invertBoolean returns the negation of a "true"/"false" string
func invertBoolean(val string) bool {
	return strings.ToLower(val) != "true"
}

// extractSlug extracts the slug from post_id
func extractSlug(postID string) string {
	parts := strings.SplitN(postID, ".", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return postID
}

// loadBatchOutput reads and parses the file containing the OpenAI batch API output
func loadBatchOutput(filename string) (map[string]SEOContent, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	seoData := make(map[string]SEOContent)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var record OutputRecord
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			log.Printf("Skipping invalid JSON line: %v", err)
			continue
		}

		parts := strings.Split(record.Response.Body.Choices[0].Message.Content, "----")
		if len(parts) != 3 {

			log.Printf("Skipping malformed contents for custom_id %s", record.CustomID)
			continue
		}
		seoData[record.CustomID] = SEOContent{
			Title:       strings.TrimSpace(parts[0]),
			Description: strings.TrimSpace(parts[1]),
			Keywords:    strings.TrimSpace(parts[2]),
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return seoData, nil
}

func main() {
	ignoreDrafts := flag.Bool("ignore-drafts", false, "Ignore posts where is_published is false")
	useSEOTitle := flag.Bool("use-seo-title", false, "Use SEO-optimized title from batch_output.jsonl instead of original title")
	flag.Parse()

	fmt.Println("Starting Substack to Hugo converter...")
	if *ignoreDrafts {
		fmt.Println("Option enabled: Draft posts will be ignored.")
	} else {
		fmt.Println("Processing all posts (including drafts).")
	}

	seoMap, err := loadBatchOutput("batch_output.jsonl")
	if err != nil {
		log.Fatalf("Error loading batch_output.jsonl: %v", err)
	}

	posts, err := readPostsCSV("posts.csv")
	if err != nil {
		log.Fatalf("Error reading posts.csv: %v", err)
	}

	hugoFolder := "hugohtml"
	if err := cleanOrCreateFolder(hugoFolder); err != nil {
		log.Fatalf("Error preparing hugohtml folder: %v", err)
	}

	successCount := 0
	errors := []string{}

	for _, post := range posts {
		if *ignoreDrafts && strings.ToLower(post.IsPublished) == "false" {
			continue
		}

		slug := extractSlug(post.PostID)
		inputPath := filepath.Join("posts", post.PostID+".html")
		outputPath := filepath.Join(hugoFolder, slug+".html")

		htmlContent, err := os.ReadFile(inputPath)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed reading %s: %v", inputPath, err))
			continue
		}

		seo := seoMap[post.PostID]
		title := post.Title
		if *useSEOTitle && seo.Title != "" {
			title = seo.Title
		}

		frontMatter := fmt.Sprintf(`+++
date = %s
draft = %t
title = '%s'
weight = 10
markup = 'text/html'
slug = '%s'
description = '%s'
keywords = '%s'
[params]
  author = 'Manuel Delgado Tenorio'
+++

`,
			post.PostDate,
			invertBoolean(post.IsPublished),
			title,
			slug,
			seo.Description,
			seo.Keywords,
		)

		finalContent := frontMatter + string(htmlContent)

		if err := os.WriteFile(outputPath, []byte(finalContent), 0644); err != nil {
			errors = append(errors, fmt.Sprintf("Failed to write %s: %v", outputPath, err))
			continue
		}
		successCount++
	}

	absPath, _ := filepath.Abs(hugoFolder)
	fmt.Println("\nSummary:")
	fmt.Printf("Successfully created %d files.\n", successCount)
	if len(errors) > 0 {
		fmt.Printf("Encountered %d errors:\n", len(errors))
		for _, e := range errors {
			fmt.Println("-", e)
		}
	}
	fmt.Printf("Hugo HTML files saved in: %s\n", absPath)
}
