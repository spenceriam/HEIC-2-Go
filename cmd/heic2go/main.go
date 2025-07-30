package main

import (
	"fmt"
	"os"
)

const (
	appName    = "HEIC-2-Go"
	appVersion = "0.1.0"
)

func main() {
	fmt.Printf("%s v%s\n", appName, appVersion)
	fmt.Println("A tool to convert HEIC images to JPG format")
	fmt.Println("\nThis is a work in progress. Check back soon for updates!")
	os.Exit(0)
}
