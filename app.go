package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"murican-to-english/pkg/converter"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx       context.Context
	converter *converter.Converter
	filePath  string // Store the path of the file being processed
}

// ServiceHandler represents a macOS service handler
type ServiceHandler interface {
	HandleService(pboard string, userData string) string
	HandleFileService(fileURL string) error
}

// Dictionary represents a mapping between words
type Dictionary map[string]string

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize the converter
	var err error
	a.converter, err = converter.NewConverter()
	if err != nil {
		fmt.Printf("Error initializing converter: %v\n", err)
	}

	// Check if the app was launched with a file path argument
	args := os.Args
	if len(args) > 1 {
		filePath := args[1]
		// Check if the file exists
		if _, err := os.Stat(filePath); err == nil {
			a.filePath = filePath
			// Process the file and exit
			err := a.ConvertFileToEnglish(filePath)
			if err != nil {
				fmt.Printf("Error converting file: %v\n", err)
			}
			// Exit the application after processing the file
			os.Exit(0)
		}
	}
}

// domReady is called when the DOM is ready
func (a *App) domReady(ctx context.Context) {
	// Try multiple approaches to center the window

	// First, try the built-in center function
	runtime.WindowCenter(ctx)

	// Then set the window position to a value that's likely to be centered on most displays
	// These values are chosen to position the window more centrally on common screen resolutions
	runtime.WindowSetPosition(ctx, 500, 300)

	// Try a different position that might work better on larger displays
	runtime.WindowSetPosition(ctx, 600, 350)

	// Finally, try the built-in center function again
	runtime.WindowCenter(ctx)
}

// ConvertToBritish converts American English text to British English
func (a *App) ConvertToBritish(text string, normaliseSmartQuotes bool) string {
	if a.converter == nil {
		return "Error: Converter not initialized"
	}
	return a.converter.ConvertToBritish(text, normaliseSmartQuotes)
}

// ConvertToAmerican converts British English text to American English
func (a *App) ConvertToAmerican(text string, normaliseSmartQuotes bool) string {
	if a.converter == nil {
		return "Error: Converter not initialized"
	}
	return a.converter.ConvertToAmerican(text, normaliseSmartQuotes)
}

// ConvertFileToEnglish converts a file's content from American to British English and saves it back
func (a *App) ConvertFileToEnglish(filePath string) error {
	// Read the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Convert the content
	convertedContent := a.ConvertToBritish(string(content), true)

	// Write the converted content back to the file
	err = ioutil.WriteFile(filePath, []byte(convertedContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

// HandleDroppedFile processes a file that was dropped onto the application
func (a *App) HandleDroppedFile(filePath string) (string, error) {
	// Store the file path
	a.filePath = filePath

	// Read the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	// Return the content for display in the UI
	return string(content), nil
}

// SaveConvertedFile saves the converted content back to the original file
func (a *App) SaveConvertedFile(content string) error {
	if a.filePath == "" {
		return fmt.Errorf("no file path set")
	}

	// Write the content back to the file
	err := ioutil.WriteFile(a.filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	// Clear the file path after saving
	a.filePath = ""

	return nil
}

// GetCurrentFilePath returns the path of the currently loaded file
func (a *App) GetCurrentFilePath() string {
	return a.filePath
}

// ClearCurrentFile clears the current file path
func (a *App) ClearCurrentFile() {
	a.filePath = ""
}

// GetAmericanToBritishDictionary returns the American to British dictionary
func (a *App) GetAmericanToBritishDictionary() Dictionary {
	if a.converter == nil {
		return Dictionary{}
	}
	return a.converter.GetAmericanToBritishDictionary()
}

// GetBritishToAmericanDictionary returns the British to American dictionary
func (a *App) GetBritishToAmericanDictionary() Dictionary {
	if a.converter == nil {
		return Dictionary{}
	}
	return a.converter.GetBritishToAmericanDictionary()
}

// HandleService processes text from the macOS service menu
func (a *App) HandleService(pboard string, userData string) string {
	// Convert the text to British English
	if a.converter == nil {
		// Initialize the converter if it's not already initialized
		var err error
		a.converter, err = converter.NewConverter()
		if err != nil {
			return "Error initializing converter: " + err.Error()
		}
	}

	// Convert the text
	return a.ConvertToBritish(pboard, true)
}

// HandleFileService processes a file from the macOS service menu
func (a *App) HandleFileService(fileURL string) error {
	// Convert the file URL to a file path
	// macOS file URLs are in the format "file:///path/to/file"
	filePath := strings.TrimPrefix(fileURL, "file://")

	// Initialize the converter if it's not already initialized
	if a.converter == nil {
		var err error
		a.converter, err = converter.NewConverter()
		if err != nil {
			return fmt.Errorf("error initializing converter: %w", err)
		}
	}

	// Convert the file
	return a.ConvertFileToEnglish(filePath)
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	// Perform any cleanup or save settings here
}
