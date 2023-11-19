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
func processTemplate(templatePath, outputPath, templateName string, config Config) {
	inputPath := filepath.Join(templatePath, templateName)
	input, err := os.ReadFile(inputPath)
	if err != nil {
		logError("Failed to read template: "+templateName, err)
		return
	}

	content := string(input)

	// Extract and remove compliance tags, and collect compliance IDs
	var complianceIDs []string
	content, complianceIDs = extractAndRemoveComplianceTags(content)

	//Ugly debug
	//fmt.Println("Extracted Compliance IDs:", complianceIDs)

	// Append compliance IDs to the Compliance section
	content = appendComplianceSection(content, complianceIDs)
	//Ugly debug
	//fmt.Println("Final Content:", content)

	// Replace config variables
	content = strings.ReplaceAll(content, "{{companyName}}", config.CompanyName)
	content = strings.ReplaceAll(content, "{{ceoName}}", config.CeoName)
	content = strings.ReplaceAll(content, "{{cisoName}}", config.CisoName)
	content = strings.ReplaceAll(content, "{{mainSecPolicyName}}", config.MainSecPolicyName)

	outputFile := filepath.Join(outputPath, templateName)
	err = os.WriteFile(outputFile, []byte(content), 0644)
	if err != nil {
		logError("Failed to write processed template: "+templateName, err)
		return
	}

	fmt.Println("Processed and saved:", outputFile)
}

// extractAndRemoveComplianceTags extracts compliance IDs and removes the tags
func extractAndRemoveComplianceTags(content string) (string, []string) {
	idMap := make(map[string]bool)
	var uniqueIDs []string

	// Regex pattern to match <!-- compliance: ... -->
	pattern := `<!-- compliance: (.*?) -->`
	re := regexp.MustCompile(pattern)

	// Find all matches and add unique IDs to the map
	matches := re.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			id := match[1]
			if _, exists := idMap[id]; !exists {
				idMap[id] = true
				uniqueIDs = append(uniqueIDs, id)
			}
		}
	}

	// Remove the tags
	content = re.ReplaceAllString(content, "")
	content = regexp.MustCompile(`<!-- /compliance: .*? -->`).ReplaceAllString(content, "")

	return content, uniqueIDs
}

// appendComplianceSection appends the compliance IDs to the Compliance section
func appendComplianceSection(content string, ids []string) string {
	if len(ids) == 0 {
		return content
	}

	// Prepare the string to be inserted
	complianceInsert := "\nCompliance controls covered by the document:\n"
	for _, id := range ids {
		complianceInsert += "- " + id + "\n"
	}

	// Define the pattern to find the Compliance Control section
	pattern := `## Compliance Control.*\n`
	re := regexp.MustCompile(pattern)

	// Find the location of the Compliance Control section
	location := re.FindStringIndex(content)
	if location != nil {
		// Insert the compliance IDs into the existing section
		before := content[:location[1]]
		after := content[location[1]:]
		return before + complianceInsert + after
	}

	// If the Compliance Control section is not found, append at the end
	return content + "\n## Compliance Control" + complianceInsert
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

func logError(message string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
}
