package terminal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/localrivet/gomcp/server"
)

// CatArgs defines the arguments for the cat tool
type CatArgs struct {
	Path string `json:"path"`
}

// CatResult defines the result structure for the cat tool
type CatResult struct {
	Content  string `json:"content"`
	FilePath string `json:"file_path"`
	Size     int64  `json:"size"`
}

// HandleCat implements the logic for the cat tool
// This handler reads and returns the content of the file at the provided path
func HandleCat(ctx *server.Context, args CatArgs) (CatResult, error) {
	ctx.Logger.Info("Handling Cat tool call")

	// Validate the path
	if args.Path == "" {
		return CatResult{}, fmt.Errorf("path cannot be empty")
	}

	// Validate path for security
	if err := validatePath(ctx, args.Path); err != nil {
		return CatResult{}, fmt.Errorf("path validation failed: %w", err)
	}

	// Clean and resolve the path
	cleanPath := filepath.Clean(args.Path)
	ctx.Logger.Info("reading file", "path", cleanPath)

	// Check if file exists
	fileInfo, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return CatResult{}, fmt.Errorf("file does not exist: %s", cleanPath)
		}
		return CatResult{}, fmt.Errorf("failed to access file: %w", err)
	}

	// Check if it's a directory
	if fileInfo.IsDir() {
		return CatResult{}, fmt.Errorf("cannot cat a directory: %s", cleanPath)
	}

	// Check file size (optional safety check for very large files)
	fileSize := fileInfo.Size()
	const maxFileSize = 10 * 1024 * 1024 // 10MB limit
	if fileSize > maxFileSize {
		return CatResult{}, fmt.Errorf("file too large to cat: %d bytes (max %d bytes)", fileSize, maxFileSize)
	}

	// Read the file content
	content, err := os.ReadFile(cleanPath)
	if err != nil {
		return CatResult{}, fmt.Errorf("failed to read file: %w", err)
	}

	// Convert to string and handle potential binary content
	contentStr := string(content)
	if containsBinaryData(content) {
		ctx.Logger.Info("File appears to contain binary data", "path", cleanPath)
		contentStr = fmt.Sprintf("[Binary file - %d bytes]", len(content))
	}

	result := CatResult{
		Content:  contentStr,
		FilePath: cleanPath,
		Size:     fileSize,
	}

	ctx.Logger.Info("Successfully read file", "path", cleanPath, "size", fileSize)
	return result, nil
}

// containsBinaryData checks if the content appears to be binary
func containsBinaryData(data []byte) bool {
	// Simple heuristic: if more than 10% of the first 1024 bytes are non-printable, consider it binary
	sampleSize := 1024
	if len(data) < sampleSize {
		sampleSize = len(data)
	}

	nonPrintableCount := 0
	for i := 0; i < sampleSize; i++ {
		b := data[i]
		// Consider control characters (except tab, newline, carriage return) as non-printable
		if b < 32 && b != 9 && b != 10 && b != 13 {
			nonPrintableCount++
		}
		// Also check for high-bit characters that might indicate binary
		if b > 126 {
			nonPrintableCount++
		}
	}

	return float64(nonPrintableCount)/float64(sampleSize) > 0.1
}

// validatePath performs basic security checks on the file path
func validatePath(ctx *server.Context, path string) error {
	// Prevent directory traversal attacks
	if strings.Contains(path, "..") {
		return fmt.Errorf("path traversal not allowed: %s", path)
	}

	// Convert to absolute path for additional checks
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Additional security checks could be added here
	// For example, checking against allowed directories from config
	ctx.Logger.Info("Validated path", "original", path, "absolute", absPath)

	return nil
}
