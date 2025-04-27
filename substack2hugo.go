package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
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

	// Map header indexes
	fieldIdx := make(map[string]int)
	for i, h := range headers {
		fieldIdx[h] = i
	}

	requiredFields := []string{"post_id", "post_date", "is_published", "title"}
	for _, field := range requiredFields {
		if _, ok := fieldIdx[field]; !ok {
			return nil, fmt.Errorf("missing required field in CSV: %s", field)
		}
	}

	var posts []Post
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV record: %w", err)
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

	dirEntries, err := os.ReadDir(folder)
	if err != nil {
		return err
	}

	for _, entry := range dirEntries {
		entryPath := filepath.Join(folder, entry.Name())
		err = os.RemoveAll(entryPath)
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

func main() {
	// Define the --ignore-drafts flag
	ignoreDrafts := flag.Bool("ignore-drafts", false, "Ignore posts where is_published is false")
	flag.Parse()

	fmt.Println("Starting Substack to Hugo converter...")
	if *ignoreDrafts {
		fmt.Println("Option enabled: Draft posts will be ignored.")
	} else {
		fmt.Println("Processing all posts (including drafts).")
	}

	posts, err := readPostsCSV("posts.csv")
	if err != nil {
		log.Fatalf("Error reading posts.csv: %v", err)
	}

	hugoFolder := "hugohtml"
	err = cleanOrCreateFolder(hugoFolder)
	if err != nil {
		log.Fatalf("Error preparing hugohtml folder: %v", err)
	}

	successCount := 0
	errors := make([]string, 0)

	for _, post := range posts {
		if *ignoreDrafts && strings.ToLower(post.IsPublished) == "false" {
			continue // Skip draft
		}

		slug := extractSlug(post.PostID)
		inputPath := filepath.Join("posts", post.PostID+".html")
		outputPath := filepath.Join(hugoFolder, slug+".html")

		content, err := os.ReadFile(inputPath)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed reading %s: %v", inputPath, err))
			continue
		}

		frontMatter := fmt.Sprintf(`+++
date = %s
draft = %t
title = '%s'
weight = 10
markup = 'text/html'
slug = '%s'
[params]
  author = 'Manuel Delgado Tenorio'
+++

`,
			post.PostDate,
			invertBoolean(post.IsPublished),
			post.Title,
			slug,
		)

		fullContent := frontMatter + string(content)

		err = os.WriteFile(outputPath, []byte(fullContent), 0644)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed writing %s: %v", outputPath, err))
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