package maestro

import (
	"archive/zip"
	"bufio"
	"encoding/json"
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
	baseURL        = "https://open.devicelab.dev"
	downloadURLFmt = baseURL + "/download/maestro-complete-reports/jars/%s.zip"
	listVersionURL = baseURL + "/api/maestro-complete-reports/jars"
	backupDir      = ".maestro/backup"
)

// ErrVersionNotSupported is returned when the Maestro version is not supported
var ErrVersionNotSupported = fmt.Errorf("version not supported")

// VersionListResponse represents the API response for supported versions
type VersionListResponse struct {
	Project  string   `json:"project"`
	Versions []string `json:"versions"`
}

type Maestro struct {
	Version string
	LibPath string
}

func Detect() (*Maestro, error) {
	version, err := getVersion()
	if err != nil {
		return nil, fmt.Errorf("maestro not found: %w", err)
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

// DownloadAndReplaceJars downloads JARs for the detected Maestro version.
// If 404, it fetches supported versions and returns ErrVersionNotSupported.
func (m *Maestro) DownloadAndReplaceJars() error {
	downloadURL := fmt.Sprintf(downloadURLFmt, m.Version)

	tempDir, err := os.MkdirTemp("", "maestro-jars-")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	zipPath := filepath.Join(tempDir, "jars.zip")

	// Download - check for 404 (version not supported)
	fmt.Println("Downloading patched JARs...")
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download jars: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// Version not supported - fetch list and return error
		versions, listErr := GetSupportedVersions()
		if listErr != nil {
			return fmt.Errorf("maestro version %s is not supported", m.Version)
		}
		return fmt.Errorf("maestro version %s is not supported. Supported versions: %v", m.Version, versions)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download jars: %s", resp.Status)
	}

	// Save zip file
	out, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, resp.Body)
	out.Close()
	if err != nil {
		return fmt.Errorf("failed to save zip file: %w", err)
	}

	// Extract
	extractPath := filepath.Join(tempDir, "extracted")
	if err := unzip(zipPath, extractPath); err != nil {
		return fmt.Errorf("failed to extract jars: %w", err)
	}

	// Find and copy JARs to lib directory (handles nested directories)
	fmt.Println("Installing JARs:")
	jarCount := 0
	err = filepath.Walk(extractPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		// Skip macOS metadata files
		if strings.Contains(path, "__MACOSX") {
			return nil
		}
		if strings.HasSuffix(info.Name(), ".jar") {
			dst := filepath.Join(m.LibPath, info.Name())
			if err := copyFile(path, dst); err != nil {
				return fmt.Errorf("failed to copy %s: %w", info.Name(), err)
			}
			fmt.Printf("  ✓ %s\n", dst)
			jarCount++
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to install jars: %w", err)
	}

	if jarCount == 0 {
		return fmt.Errorf("no JAR files found in downloaded archive")
	}

	return nil
}

// GetSupportedVersions fetches the list of supported Maestro versions from the API
func GetSupportedVersions() ([]string, error) {
	resp, err := http.Get(listVersionURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch versions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch versions: %s", resp.Status)
	}

	var result VersionListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse versions: %w", err)
	}

	return result.Versions, nil
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
