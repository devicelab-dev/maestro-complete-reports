package maestro

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	// TODO: Replace with actual server endpoint
	jarDownloadURL = "PLACEHOLDER_URL"
	backupDir      = ".maestro/backup"
)

var supportedVersions = []string{"2.0.9", "2.0.10"}

type Maestro struct {
	Version string
	LibPath string
}

func Detect() (*Maestro, error) {
	version, err := getVersion()
	if err != nil {
		return nil, fmt.Errorf("maestro not found: %w", err)
	}

	if !isSupportedVersion(version) {
		return nil, fmt.Errorf("unsupported maestro version: %s (supported: %v)", version, supportedVersions)
	}

	libPath, err := getLibPath()
	if err != nil {
		return nil, fmt.Errorf("could not find maestro lib directory: %w", err)
	}

	return &Maestro{
		Version: version,
		LibPath: libPath,
	}, nil
}

func getVersion() (string, error) {
	cmd := exec.Command("maestro", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	patterns := []string{
		`cli version:\s*(\d+\.\d+\.\d+)`,
		`version:\s*(\d+\.\d+\.\d+)`,
		`CLI\s+(\d+\.\d+\.\d+)`,
		`(\d+\.\d+\.\d+)`,
	}

	outputStr := string(output)
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(outputStr)
		if len(matches) > 1 {
			return matches[1], nil
		}
	}

	return "", fmt.Errorf("could not parse version from output: %s", outputStr)
}

func isSupportedVersion(version string) bool {
	for _, v := range supportedVersions {
		if v == version {
			return true
		}
	}
	return false
}

func getLibPath() (string, error) {
	// Strategy 1: Parse maestro script for CLASSPATH
	if path, err := getLibPathFromScript(); err == nil {
		return path, nil
	}

	// Strategy 2: Check sibling lib directory
	if path, err := getLibPathFromSibling(); err == nil {
		return path, nil
	}

	// Strategy 3: Fallback to ~/.maestro/lib
	if path, err := getLibPathFromHome(); err == nil {
		return path, nil
	}

	return "", fmt.Errorf("could not find lib directory")
}

func getLibPathFromScript() (string, error) {
	cmd := exec.Command("which", "maestro")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	scriptPath := strings.TrimSpace(string(output))
	file, err := os.Open(scriptPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`(.*/lib/)`)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "CLASSPATH") || strings.Contains(line, "/lib/") {
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				libPath := matches[1]
				if _, err := os.Stat(libPath); err == nil {
					return libPath, nil
				}
			}
		}
	}

	return "", fmt.Errorf("no lib path found in script")
}

func getLibPathFromSibling() (string, error) {
	cmd := exec.Command("which", "maestro")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	scriptPath := strings.TrimSpace(string(output))
	binDir := filepath.Dir(scriptPath)
	parentDir := filepath.Dir(binDir)
	libPath := filepath.Join(parentDir, "lib")

	if _, err := os.Stat(libPath); err == nil {
		return libPath, nil
	}

	return "", fmt.Errorf("no sibling lib directory found")
}

func getLibPathFromHome() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	libPath := filepath.Join(homeDir, ".maestro", "lib")
	if _, err := os.Stat(libPath); err == nil {
		return libPath, nil
	}

	return "", fmt.Errorf("no lib directory found in home")
}

func (m *Maestro) BackupJars() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	backupPath := filepath.Join(homeDir, backupDir)
	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	entries, err := os.ReadDir(m.LibPath)
	if err != nil {
		return "", fmt.Errorf("failed to read lib directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, "maestro") && strings.HasSuffix(name, ".jar") {
			src := filepath.Join(m.LibPath, name)
			dst := filepath.Join(backupPath, name)
			if err := copyFile(src, dst); err != nil {
				return "", fmt.Errorf("failed to backup %s: %w", name, err)
			}
		}
	}

	return backupPath, nil
}

func (m *Maestro) DownloadAndReplaceJars(downloadURL string) error {
	tempDir, err := os.MkdirTemp("", "maestro-jars-")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	zipPath := filepath.Join(tempDir, "jars.zip")

	// Download
	if err := downloadFile(downloadURL, zipPath); err != nil {
		return fmt.Errorf("failed to download jars: %w", err)
	}

	// Extract
	extractPath := filepath.Join(tempDir, "extracted")
	if err := unzip(zipPath, extractPath); err != nil {
		return fmt.Errorf("failed to extract jars: %w", err)
	}

	// Copy JARs to lib directory
	entries, err := os.ReadDir(extractPath)
	if err != nil {
		return fmt.Errorf("failed to read extracted directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".jar") {
			src := filepath.Join(extractPath, name)
			dst := filepath.Join(m.LibPath, name)
			if err := copyFile(src, dst); err != nil {
				return fmt.Errorf("failed to copy %s: %w", name, err)
			}
		}
	}

	return nil
}

func (m *Maestro) RestoreJars() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	backupPath := filepath.Join(homeDir, backupDir)
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("no backup found at %s", backupPath)
	}

	entries, err := os.ReadDir(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".jar") {
			src := filepath.Join(backupPath, name)
			dst := filepath.Join(m.LibPath, name)
			if err := copyFile(src, dst); err != nil {
				return fmt.Errorf("failed to restore %s: %w", name, err)
			}
		}
	}

	return nil
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		// Security: prevent zip slip
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0755)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
