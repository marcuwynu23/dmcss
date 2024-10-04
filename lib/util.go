package lib
import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Utility functions like readAndProcessFile, transformTokensToCSS, etc. go here
// These functions are shared between compile, generate, and script subcommands

// Token structure to store device query data
type Token struct {
	Name   string
	Width  string
	Height string
	Body   string
}

func tokenizeCustomSyntax(input string, deviceDir string) ([]Token, string, error) {
	regex := regexp.MustCompile(`\$\$device\s*\(\s*name:\s*"([a-zA-Z0-9-]+)",\s*width:\s*([a-zA-Z0-9%]+),\s*height:\s*([a-zA-Z0-9%]+)\)\s*`)
	matches := regex.FindAllStringSubmatch(input, -1)

	var tokens []Token
	processedInput := input

	for _, match := range matches {
		deviceName := match[1]
		deviceFile := filepath.Join(deviceDir, deviceName+".dmcss")

		// Read the device-specific CSS from the corresponding file
		deviceContent, err := ioutil.ReadFile(deviceFile)
		if err != nil {
			return nil, "", fmt.Errorf("error reading device file %s: %v", deviceFile, err)
		}

		// Add token
		tokens = append(tokens, Token{
			Name:   deviceName,
			Width:  match[2],
			Height: match[3],
			Body:   strings.TrimSpace(string(deviceContent)),
		})

		// Remove $$device line from the input
		processedInput = strings.Replace(processedInput, match[0], "", 1)
	}

	return tokens, processedInput, nil
}

func transformTokensToCSS(tokens []Token) string {
	var cssBuilder strings.Builder

	for _, token := range tokens {
		cssBuilder.WriteString(fmt.Sprintf("@media (width: %s) and (height: %s) {\n", token.Width, token.Height))
		bodyLines := strings.Split(token.Body, "\n")
		for _, line := range bodyLines {
			trimmedLine := strings.TrimSpace(line)
			if strings.HasSuffix(trimmedLine, "{") {
				cssBuilder.WriteString(fmt.Sprintf("  .device-%s %s {\n", token.Name, strings.Split(trimmedLine, "{")[0]))
			} else if trimmedLine == "}" {
				cssBuilder.WriteString("  }\n")
			} else {
				cssBuilder.WriteString(fmt.Sprintf("    %s\n", trimmedLine))
			}
		}
		cssBuilder.WriteString("}\n\n")
	}

	return cssBuilder.String()
}

func readAndProcessFile(filePath string, deviceDir string) (string, []Token, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", nil, err
	}

	processedContent, err := processImports(string(content), filepath.Dir(filePath))
	if err != nil {
		return "", nil, err
	}

	tokens, processedInput, err := tokenizeCustomSyntax(processedContent, deviceDir)
	if err != nil {
		return "", nil, err
	}

	return processedInput, tokens, nil
}

func processImports(input string, baseDir string) (string, error) {
	importRegex := regexp.MustCompile(`@import\s*"([^"]+)"\s*;`)
	matches := importRegex.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		importPath := match[1]
		fullImportPath := filepath.Join(baseDir, importPath)
		importContent, err := ioutil.ReadFile(fullImportPath)
		if err != nil {
			return "", fmt.Errorf("error importing file %s: %v", importPath, err)
		}
		importedContent, err := processImports(string(importContent), baseDir)
		if err != nil {
			return "", err
		}
		input = strings.Replace(input, match[0], importedContent, 1)
	}

	return input, nil
}

func writeOutputFile(filePath string, data string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory %s: %v", dir, err)
	}

	return ioutil.WriteFile(filePath, []byte(data), 0644)
}

func generateDeviceFile(deviceName, width, height, deviceDir string) error {
	deviceFilePath := filepath.Join(deviceDir, deviceName+".dmcss")
	deviceContent := fmt.Sprintf("/* Styles for %s */\nbody {\n    /* Add styles here */\n}\n", deviceName)
	err := ioutil.WriteFile(deviceFilePath, []byte(deviceContent), 0644)
	if err != nil {
		return fmt.Errorf("error creating device file %s: %v", deviceFilePath, err)
	}

	mainFilePath := filepath.Join(deviceDir, "../index.dmcss")
	deviceLine := fmt.Sprintf("\n$$device(name: \"%s\", width: %spx, height: %spx)\n", deviceName, width, height)
	err = appendToFile(mainFilePath, deviceLine)
	if err != nil {
		return fmt.Errorf("error appending to index.dmcss: %v", err)
	}

	fmt.Printf("Device %s successfully generated and appended to index.dmcss.\n", deviceName)
	return nil
}

func appendToFile(filePath, text string) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		return err
	}
	return nil
}


// Function to extract device names from index.dmcss for script generation
func extractDeviceNames(input string) ([]string, error) {
	// Regex to match custom $$device queries like: $$device(name: "iphone-xr", width: 414px, height: 896px)
	regex := regexp.MustCompile(`\$\$device\s*\(\s*name:\s*"([a-zA-Z0-9-]+)",\s*width:\s*([a-zA-Z0-9%]+),\s*height:\s*([a-zA-Z0-9%]+)\)\s*`)
	matches := regex.FindAllStringSubmatch(input, -1)

	var deviceNames []string
	// Loop over the matches to extract the device names
	for _, match := range matches {
		deviceNames = append(deviceNames, match[1])
	}

	if len(deviceNames) == 0 {
		return nil, fmt.Errorf("no device names found in the input")
	}

	return deviceNames, nil
}


func generateScriptFile(deviceNames []string, outputFilePath string) error {
	var jsBuilder strings.Builder

	// Start building the script content
	jsBuilder.WriteString("document.addEventListener('DOMContentLoaded', function() {\n")
	jsBuilder.WriteString("  const body = document.body;\n")
	jsBuilder.WriteString("  let classNames = [];\n")

	// Add each device name to the classNames array
	for _, device := range deviceNames {
		jsBuilder.WriteString(fmt.Sprintf("  classNames.push('device-%s');\n", device))
	}

	// Prepend the device classes to the existing body className
	jsBuilder.WriteString("  body.className = classNames.join(' ') + ' ' + body.className;\n")
	jsBuilder.WriteString("});\n")

	// Write the generated JavaScript to the file
	return ioutil.WriteFile(outputFilePath, []byte(jsBuilder.String()), 0644)
}

// Subcommand to generate script.js
func generateScriptCmd(inputFilePath, outputFilePath string) error {
	// Read the index.dmcss file
	content, err := ioutil.ReadFile(inputFilePath)
	if err != nil {
		return fmt.Errorf("error reading index.dmcss file: %v", err)
	}

	// Extract device names from the index.dmcss file
	deviceNames, err := extractDeviceNames(string(content))
	if err != nil {
		return fmt.Errorf("error extracting device names: %v", err)
	}

	// Generate the script.js file
	err = generateScriptFile(deviceNames, outputFilePath)
	if err != nil {
		return fmt.Errorf("error generating script.js: %v", err)
	}

	fmt.Printf("JavaScript file successfully generated at %s\n", outputFilePath)
	return nil
}