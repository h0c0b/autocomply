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
	CompanyName       string `json:"companyName"`
	CeoName           string `json:"ceoName"`
	CisoName          string `json:"cisoName"`
	MainSecPolicyName string `json:"mainSecPolicyName"`
}

func main() {
	// Check command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: autocomply <command>")
		fmt.Println("Commands:\n  build\n  report")
		return
	}

	command := os.Args[1]

	switch command {
	case "build":
		buildProject()
	case "report":
		generateReport()
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}

func buildProject() {
	cwd, err := os.Getwd()
	if err != nil {
		logError("Failed to get current directory", err)
		return
	}

	config, err := loadConfig(filepath.Join(cwd, "config.json"))
	if err != nil {
		logError("Failed to load configuration", err)
		return
	}

	outputPath := filepath.Join(cwd, "output")
	err = os.MkdirAll(outputPath, 0755)
	if err != nil {
		logError("Failed to create output directory", err)
		return
	}

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

func generateReport() {
	// Implement report generation logic here
	fmt.Println("Report generation not implemented yet.")
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
	content = strings.ReplaceAll(content, "{{mainSecPolicyName}}", config.MainSecPolicyName)

	// Remove 'compliance' tags
	content = removeComplianceTags(content)

	outputFile := filepath.Join(outputPath, templateName)
	err = os.WriteFile(outputFile, []byte(content), 0644)
	if err != nil {
		logError("Failed to write processed template: "+templateName, err)
		return
	}

	fmt.Println("Processed and saved:", outputFile)
}

// removeComplianceTags removes compliance HTML comment tags
func removeComplianceTags(content string) string {
	// Regex pattern to match <!-- compliance: ... --> and <!-- /compliance: ... --> tags
	pattern := `<!-- compliance: .*? -->|<!-- /compliance: .*? -->`
	re := regexp.MustCompile(pattern)

	// Replace all occurrences with an empty string
	return re.ReplaceAllString(content, "")
}

func logError(message string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
}
