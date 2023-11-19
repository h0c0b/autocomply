package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Config represents the structure of your JSON config file
type Config struct {
	CompanyName string `json:"companyName"`
	CeoName     string `json:"ceoName"`
	CisoName    string `json:"cisoName"`
}

func main() {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		logError("Failed to get current directory", err)
		return
	}

	// Load the configuration
	config, err := loadConfig(filepath.Join(cwd, "config.json"))
	if err != nil {
		logError("Failed to load configuration", err)
		return
	}

	// Ensure output directory exists
	outputPath := filepath.Join(cwd, "output")
	err = os.MkdirAll(outputPath, 0755)
	if err != nil {
		logError("Failed to create output directory", err)
		return
	}

	// Process each Markdown file
	templatePath := filepath.Join(cwd, "templates")
	templates, err := os.ReadDir(templatePath)
	if err != nil {
		logError("Failed to read templates directory", err)
		return
	}

	for _, template := range templates {
		processTemplate(templatePath, outputPath, template.Name(), config)
	}
}

func loadConfig(filePath string) (Config, error) {
	var config Config
	file, err := os.ReadFile(filePath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(file, &config)
	return config, err
}

func processTemplate(templatePath, outputPath, templateName string, config Config) {
	inputPath := filepath.Join(templatePath, templateName)
	input, err := os.ReadFile(inputPath)
	if err != nil {
		logError("Failed to read template: "+templateName, err)
		return
	}

	content := string(input)

	// Replace config variables
	content = strings.ReplaceAll(content, "{{companyName}}", config.CompanyName)
	content = strings.ReplaceAll(content, "{{ceoName}}", config.CeoName)
	content = strings.ReplaceAll(content, "{{cisoName}}", config.CisoName)

	// Remove 'implements' variables
	content = removeImplementsVariables(content)

	outputFile := filepath.Join(outputPath, templateName)
	err = os.WriteFile(outputFile, []byte(content), 0644)
	if err != nil {
		logError("Failed to write processed template: "+templateName, err)
		return
	}

	fmt.Println("Processed and saved:", outputFile)
}

// removeImplementsVariables removes placeholders that start with 'implements'
func removeImplementsVariables(content string) string {
	// Regex pattern to match {{implements: ...}} variables
	pattern := `{{implements: [^}]+}}`
	re := regexp.MustCompile(pattern)

	// Replace all occurrences with an empty string
	return re.ReplaceAllString(content, "")
}

func logError(message string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
}
