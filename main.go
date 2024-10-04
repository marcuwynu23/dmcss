package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/marcuwynu23/dmcss/lib" // Import the lib package
)

func main() {
	// Create a flag set for subcommands
	compileCmd := flag.NewFlagSet("compile", flag.ExitOnError)
	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)

	// Flags for the generate command
	deviceName := generateCmd.String("device", "", "Name of the device")
	width := generateCmd.String("width", "", "Device width in px")
	height := generateCmd.String("height", "", "Device height in px")

	// Ensure we have at least one argument
	if len(os.Args) < 2 {
		fmt.Println("Expected 'compile' or 'generate' subcommands")
		os.Exit(1)
	}

	// Parse the subcommands
	switch os.Args[1] {
	case "compile":
		compileCmd.Parse(os.Args[2:])
		lib.Compile(compileCmd) // Call the Compile function from the lib package
	case "generate":
		generateCmd.Parse(os.Args[2:])
		if *deviceName == "" || *width == "" || *height == "" {
			fmt.Println("Usage: dmcss generate --device <name> --width <px> --height <px>")
			os.Exit(1)
		}
		lib.GenerateDevice(*deviceName, *width, *height) // Call the GenerateDevice function from the lib package
	case "generate-script":
		lib.GenerateScript() // Call the GenerateScript function from the lib package
	default:
		fmt.Println("Expected 'compile', 'generate', or 'generate-script' subcommands")
		os.Exit(1)
	}
}
