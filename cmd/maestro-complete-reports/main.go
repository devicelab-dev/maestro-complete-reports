package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/devicelab-dev/maestro-complete-reports/internal/maestro"
)

const (
	// TODO: Replace with actual server endpoint that returns R2 temp URL
	serverEndpoint = "PLACEHOLDER_SERVER_URL"
)

type DownloadResponse struct {
	URL string `json:"url"`
}

func main() {
	setupCmd := flag.NewFlagSet("setup", flag.ExitOnError)
	restoreCmd := flag.NewFlagSet("restore", flag.ExitOnError)

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "setup":
		setupCmd.Parse(os.Args[2:])
		runSetup()
	case "restore":
		restoreCmd.Parse(os.Args[2:])
		runRestore()
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: maestro-complete-reports <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  setup    Detect Maestro, backup JARs, download and replace with patched JARs")
	fmt.Println("  restore  Restore original JARs from backup")
}

func runSetup() {
	fmt.Println("Detecting Maestro installation...")

	m, err := maestro.Detect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found Maestro %s\n", m.Version)
	fmt.Printf("Lib path: %s\n", m.LibPath)

	fmt.Println("Backing up original JARs...")
	backupPath, err := m.BackupJars()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error backing up JARs: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Backup created at: %s\n", backupPath)

	fmt.Println("Fetching download URL from server...")
	downloadURL, err := getDownloadURL()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching download URL: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Downloading and replacing JARs...")
	if err := m.DownloadAndReplaceJars(downloadURL); err != nil {
		fmt.Fprintf(os.Stderr, "Error replacing JARs: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Setup complete!")
}

func runRestore() {
	fmt.Println("Detecting Maestro installation...")

	m, err := maestro.Detect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found Maestro %s\n", m.Version)
	fmt.Printf("Lib path: %s\n", m.LibPath)

	fmt.Println("Restoring original JARs from backup...")
	if err := m.RestoreJars(); err != nil {
		fmt.Fprintf(os.Stderr, "Error restoring JARs: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Restore complete!")
}

func getDownloadURL() (string, error) {
	resp, err := http.Get(serverEndpoint)
	if err != nil {
		return "", fmt.Errorf("failed to contact server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned status: %s", resp.Status)
	}

	var result DownloadResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if result.URL == "" {
		return "", fmt.Errorf("server returned empty URL")
	}

	return result.URL, nil
}
