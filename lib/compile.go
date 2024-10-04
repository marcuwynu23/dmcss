package lib

import (
	"fmt"
	"os"
	"flag"
)

func Compile(compileCmd *flag.FlagSet) {
	// Default file paths
	inputFilePath := "dmcss/index.dmcss"
	outputFilePath := "output/output.css" // Default output file

	// If an argument is passed, use it for the output file path
	if compileCmd.NArg() > 0 {
		outputFilePath = compileCmd.Arg(0)
	}

	// Read the input .dmcss file, process imports, and tokenize device syntax
	input, tokens, err := readAndProcessFile(inputFilePath, "dmcss/devices")
	if err != nil {
		fmt.Println("Error reading input file:", err)
		os.Exit(1)
	}

	// Transform tokens to standard CSS
	outputCSS := transformTokensToCSS(tokens)

	// Combine the processed input (which includes the @import contents) with the transformed tokens
	finalOutput := input + "\n" + outputCSS

	// Write the output CSS file
	err = writeOutputFile(outputFilePath, finalOutput)
	if err != nil {
		fmt.Println("Error writing output file:", err)
		os.Exit(1)
	}

	fmt.Printf("Transpilation successful! Output written to %s\n", outputFilePath)
}
