package dropbox

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/localrivet/gomcp/server"
)

type FilesDownloadArgs struct {
	Path string `json:"path"`
}

type FileLockInfo struct {
	Created        string `json:"created"`
	IsLockholder   bool   `json:"is_lockholder"`
	LockholderName string `json:"lockholder_name"`
}

type PropertyField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PropertyGroup struct {
	Fields     []PropertyField `json:"fields"`
	TemplateID string          `json:"template_id"`
}

type SharingInfo struct {
	ModifiedBy           string `json:"modified_by"`
	ParentSharedFolderID string `json:"parent_shared_folder_id"`
	ReadOnly             bool   `json:"read_only"`
}

type DropboxFileMetadata struct {
	ClientModified           string          `json:"client_modified"`
	ContentHash              string          `json:"content_hash"`
	FileLockInfo             FileLockInfo    `json:"file_lock_info"`
	HasExplicitSharedMembers bool            `json:"has_explicit_shared_members"`
	ID                       string          `json:"id"`
	IsDownloadable           bool            `json:"is_downloadable"`
	Name                     string          `json:"name"`
	PathDisplay              string          `json:"path_display"`
	PathLower                string          `json:"path_lower"`
	PropertyGroups           []PropertyGroup `json:"property_groups"`
	Rev                      string          `json:"rev"`
	ServerModified           string          `json:"server_modified"`
	SharingInfo              SharingInfo     `json:"sharing_info"`
	Size                     int64           `json:"size"`
}

// HandleFilesDownload implements the logic the files.download tool
// This handler downloads the file at the provided FilesDownloadArgs.Path
func HandleFilesDownload(ctx *server.Context, args FilesDownloadArgs) (DropboxFileMetadata, error) {
	ctx.Logger.Info("Handling FilesDownload tool call")

	// Get API key
	apiKey := os.Getenv("DROPBOX_API_KEY")
	if apiKey == "" {
		ctx.Logger.Info("$DROPBOX_API_KEY not set")
		return DropboxFileMetadata{}, fmt.Errorf("$DROPBOX_API_KEY not set, unable to download file")
	}

	// Create the request
	req, err := createDownloadRequest(ctx, args, apiKey)
	if err != nil {
		return DropboxFileMetadata{}, fmt.Errorf("failed to create download request: %w", err)
	}

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return DropboxFileMetadata{}, fmt.Errorf("download http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := handleFailedHttpReq(resp)
		return DropboxFileMetadata{}, err
	}

	// Parse metadata from response header
	metadata, err := parseMetadataFromResponse(resp)
	if err != nil {
		return DropboxFileMetadata{}, fmt.Errorf("failed to parse metadata: %w", err)
	}

	// Read file content
	fileContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return DropboxFileMetadata{}, fmt.Errorf("failed to read file content: %w", err)
	}

	// Save file to Desktop/wip folder
	err = saveFileToDesktop(ctx, metadata.Name, fileContent)
	if err != nil {
		return DropboxFileMetadata{}, fmt.Errorf("failed to save file to Desktop: %w", err)
	}

	ctx.Logger.Info("Successfully downloaded and saved file", "path", args.Path, "size", metadata.Size, "saved_to", "Desktop/wip")
	return metadata, nil
}

// createDownloadRequest creates the HTTP request for downloading a file
func createDownloadRequest(ctx *server.Context, args FilesDownloadArgs, apiKey string) (*http.Request, error) {
	// Validate path
	if args.Path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	ctx.Logger.Info("downloading file", "path", args.Path)

	// Create the request body (empty for download)
	req, err := http.NewRequest("POST", "https://content.dropboxapi.com/2/files/download", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "text/plain")

	// Set the API argument in the header
	apiArg := map[string]string{
		"path": args.Path,
	}
	apiArgJSON, err := json.Marshal(apiArg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal API argument: %w", err)
	}
	req.Header.Set("Dropbox-API-Arg", string(apiArgJSON))

	return req, nil
}

// parseMetadataFromResponse extracts file metadata from the Dropbox-API-Result header
func parseMetadataFromResponse(resp *http.Response) (DropboxFileMetadata, error) {
	var metadata DropboxFileMetadata

	// Get the metadata from the response header
	apiResult := resp.Header.Get("Dropbox-API-Result")
	if apiResult == "" {
		return metadata, fmt.Errorf("missing Dropbox-API-Result header")
	}

	// Unmarshal the JSON metadata
	err := json.Unmarshal([]byte(apiResult), &metadata)
	if err != nil {
		return metadata, fmt.Errorf("failed to unmarshal metadata JSON: %w", err)
	}

	return metadata, nil
}

// saveFileToDesktop saves the file content to the Desktop/wip folder
func saveFileToDesktop(ctx *server.Context, filename string, content []byte) error {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Create Desktop/wip path
	wipDir := filepath.Join(homeDir, "Desktop", "wip")

	// Create the wip directory if it doesn't exist
	err = os.MkdirAll(wipDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create wip directory: %w", err)
	}

	// Create the full file path
	filePath := filepath.Join(wipDir, filename)

	// Write the file
	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	ctx.Logger.Info("File saved successfully", "file_path", filePath)
	return nil
}
