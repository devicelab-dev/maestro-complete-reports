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

func printBanner() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║  Maestro Complete Reports - by DeviceLab.dev                 ║")
	fmt.Println("║  Stop renting devices you already own.                       ║")
	fmt.Println("║  Build your own distributed lab: https://devicelab.dev       ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func printPromo() {
	fmt.Println()
	fmt.Println("Made with ❤️  by engineers who believe quality mobile testing")
	fmt.Println("shouldn't require enterprise budgets.")
	fmt.Println()
	fmt.Println("Try DeviceLab free: https://devicelab.dev")
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

	fmt.Println("Downloading and replacing JARs...")
	if err := m.DownloadAndReplaceJars(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Setup complete!")
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
