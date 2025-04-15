package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Constants for output file names
const (
	logFileName     = "log.txt"
	existFileName   = "ip_exist.txt"   // Contains targets that responded successfully at least once
	invalidFileName = "ip_invalid.txt" // Contains targets that failed the *initial* scan
)

// appendToFile appends a line to a file safely
func appendToFile(path, line string) {
	// Use O_APPEND|O_CREATE|O_WRONLY for safe concurrent appends
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%sAppend Error (%s): %v%s\n", ColorError, path, err, ColorReset)
		return
	}
	defer f.Close()
	// WriteString is generally safe for concurrent appends when file is opened with O_APPEND
	if _, err := f.WriteString(line + "\n"); err != nil {
		fmt.Fprintf(os.Stderr, "\n%sWrite Error (%s): %v%s\n", ColorError, path, err, ColorReset)
	}
}

// createOutputStructure prepares the output directory and files
func createOutputStructure(base string) error {
	err := os.MkdirAll(base, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create base output directory %s: %w", base, err)
	}

	// Pre-create category directories and status files
	for code := range statusCodes { // Assumes statusCodes map is populated globally
		categoryDigit := code / 100
		categoryName, catOk := statusCategories[categoryDigit] // Assumes statusCategories is populated globally
		if !catOk {
			continue
		}
		dir := filepath.Join(base, categoryName)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sWarning: failed to create category directory %s: %v%s\n", ColorWarning, dir, err, ColorReset)
			continue
		}
		filePath := filepath.Join(dir, fmt.Sprintf("%d.txt", code))
		f, err := os.Create(filePath) // Creates or truncates
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sWarning: failed to create status file %s: %v%s\n", ColorWarning, filePath, err, ColorReset)
		}
		if f != nil {
			f.Close()
		}
	}

	// Pre-create auxiliary files using constants
	extras := []string{existFileName, invalidFileName, logFileName}
	for _, name := range extras {
		filePath := filepath.Join(base, name)
		f, err := os.Create(filePath) // Creates or truncates
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sWarning: failed to create extra file %s: %v%s\n", ColorWarning, filePath, err, ColorReset)
		}
		if f != nil {
			f.Close()
		}
	}
	return nil
}
