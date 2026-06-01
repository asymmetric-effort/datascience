package main

import (
	"fmt"
	"os"
)

const version = "0.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "version", "--version", "-v":
		fmt.Printf("pgmgo %s\n", version)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: pgmgo [command] [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  version    Print version information")
	fmt.Println("  help       Show this help message")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -h, --help       Show help")
	fmt.Println("  -v, --version    Show version")
}
