package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/devicelab-dev/maestro-complete-reports/internal/maestro"
)

func main() {
	printBanner()

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

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBold   = "\033[1m"
)

func printBanner() {
	fmt.Println()
	fmt.Println(colorCyan + "╔══════════════════════════════════════════════════════════════╗" + colorReset)
	fmt.Println(colorCyan + "║" + colorReset + colorBold + "  Maestro Complete Reports" + colorReset + " - by " + colorGreen + "DeviceLab.dev" + colorReset + colorCyan + "                 ║" + colorReset)
	fmt.Println(colorCyan + "║" + colorReset + colorYellow + "  Stop renting devices you already own.                       " + colorReset + colorCyan + "║" + colorReset)
	fmt.Println(colorCyan + "║" + colorReset + "  Build your own distributed lab: " + colorGreen + "https://devicelab.dev" + colorReset + colorCyan + "       ║" + colorReset)
	fmt.Println(colorCyan + "╚══════════════════════════════════════════════════════════════╝" + colorReset)
	fmt.Println()
}

func printPromo() {
	fmt.Println()
	fmt.Println("Made with " + colorRed + "❤️" + colorReset + "  by engineers who believe quality mobile testing shouldn't require enterprise budgets.")
	fmt.Println("Try DeviceLab free: " + colorGreen + "https://devicelab.dev" + colorReset)
	fmt.Println()
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

	fmt.Printf("Found Maestro %s at %s\n", m.Version, m.LibPath)

	fmt.Println("Backing up original JARs...")
	if _, err := m.BackupJars(); err != nil {
		fmt.Fprintf(os.Stderr, "Error backing up JARs: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Backup complete")

	if err := m.DownloadAndReplaceJars(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println(colorGreen + "Setup complete!" + colorReset)
	printPromo()
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
	printPromo()
}
