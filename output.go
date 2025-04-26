package main

import (
	"fmt"
	"os" // Ensure os package is imported
	"path/filepath"
	"sync"
)

// Constants for output file names
const (
	logFileName            = "log.txt"
	existFileName          = "ip_exist.txt"
	invalidFileName        = "ip_invalid.txt"
	corsVulnerableFileName = "cors_vulnerable.txt"
	unknownStatusFileName  = "unknown_status.txt"
)

// Mutex to protect file writing operations across goroutines
var fileWriteMutex sync.Mutex

// appendToFile appends a line to a file safely using a mutex
func appendToFile(path, line string) {
	fileWriteMutex.Lock()
	defer fileWriteMutex.Unlock()

	// Use O_APPEND|O_CREATE|O_WRONLY with os prefix
	// Corrected line: Use os.O_APPEND, os.O_CREATE, os.O_WRONLY
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // <-- Fixed: Added os. prefix
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%sAppend Error (%s): %v%s\n", ColorError, filepath.Base(path), err, ColorReset)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(line + "\n"); err != nil {
		fmt.Fprintf(os.Stderr, "\n%sWrite Error (%s): %v%s\n", ColorError, filepath.Base(path), err, ColorReset)
	}
}

// createOutputStructure prepares the output directory and ensures all required files exist
func createOutputStructure(base string) error {
	err := os.MkdirAll(base, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create base output directory %s: %w", base, err)
	}

	// Pre-create category directories based on statusCategories map
	for _, categoryName := range statusCategories {
		dir := filepath.Join(base, categoryName)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sWarning: failed to create category directory %s: %v%s\n", ColorWarning, dir, err, ColorReset)
		}
	}
	// Also create directory for unknown categories if needed
	unknownCatDir := filepath.Join(base, "unknown_category")
	os.MkdirAll(unknownCatDir, os.ModePerm)

	// Pre-create auxiliary files using constants
	extras := []string{
		existFileName,
		invalidFileName,
		logFileName,
		corsVulnerableFileName,
		unknownStatusFileName,
	}
	for _, name := range extras {
		filePath := filepath.Join(base, name)
		// Create file if it doesn't exist, truncate if it does
		f, err := os.Create(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sWarning: failed to create/truncate auxiliary file %s: %v%s\n", ColorWarning, filePath, err, ColorReset)
		}
		if f != nil {
			f.Close()
		}
	}
	return nil
}
