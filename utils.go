package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
)

func SaveUrlBodyToFile(url string, file_path string) error {
	file, err := os.Create(file_path)
	if err != nil {
		fmt.Println("failed to cerate file - ", err, "at loc- ", file_path)
		return err
	}
	defer file.Close()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("failed to get url - ", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != http.StatusOK/100 {
		fmt.Println("return code",resp.StatusCode)
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// 4. Stream the body to the file
	_, err = io.Copy(file, resp.Body)
	return err
}


func GetFilesAsJSON(dirPath string) (string, error) {
	// Check if directory exists
	info, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("directory does not exist")
		}
		return "", err
	}

	// Ensure it's actually a directory
	if !info.IsDir() {
		return "", fmt.Errorf("path is not a directory")
	}

	// Read directory contents
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return "", err
	}

	// Collect file names
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	slices.Reverse(files)

	// Convert to JSON
	jsonData, err := json.Marshal(files)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func toggle(b bool) bool {
	if b {
		return false
	}
	return true
}
