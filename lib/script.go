package lib
import (
	"fmt"
	"os"
)

func GenerateScript() {
	inputFilePath := "dmcss/index.dmcss"
	outputFilePath := "output/script.js"

	// Generate the script.js file
	err := generateScriptCmd(inputFilePath, outputFilePath)
	if err != nil {
		fmt.Println("Error generating script.js:", err)
		os.Exit(1)
	}

	fmt.Printf("JavaScript file successfully generated at %s\n", outputFilePath)
}
