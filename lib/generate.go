package lib

import (
	"fmt"
	"os"
)

func GenerateDevice(deviceName string, width string, height string) {
	// Generate the device file and append to index.dmcss
	err := generateDeviceFile(deviceName, width, height, "dmcss/devices")
	if err != nil {
		fmt.Println("Error generating device:", err)
		os.Exit(1)
	}
}
