package dropbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/localrivet/gomcp/server"
)

type ListDropboxFoldersArgs struct {
	Path string `json:"path"`
}

type DropboxFolders []DropboxFolder
type DropboxFolder struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	PathDisplay    string `json:"path_display"`
	Tag            string `json:".tag"`
	SharedFolderID string `json:"shared_folder_id"`
}

const DROPBOX_API_URL = "https://api.dropboxapi.com/2/files/list_folder"

// HandleListDropBoxFolders implements the logic the list_dropbox_folders tool
// This handler provides a listing of all folders and their metadata at
// the provided path (ListDropboxFoldersArgs.Path).
func HandleListDropboxFolders(ctx *server.Context, args ListDropboxFoldersArgs) (DropboxFolders, error) {
	ctx.Logger.Info("Handling ListDropboxFolders tool call")

	// Get API key and print first two letters
	apiKey := os.Getenv("DROPBOX_API_KEY")
	if apiKey == "" {
		ctx.Logger.Info("$DROPBOX_API_KEY not set")
		return nil, fmt.Errorf("$DROPBOX_API_KEY not set, unable to retrieve dropbox folders")
	}
	if len(apiKey) >= 2 {
		ctx.Logger.Info("First two letters of API key: " + apiKey[:2])
	}
	// Make HTTP request to Dropbox API
	requestBody := map[string]any{
		"include_deleted":                     false,
		"include_has_explicit_shared_members": false,
		"include_media_info":                  true,
		"include_mounted_folders":             true,
		"include_non_downloadable_files":      true,
		"path":                                args.Path,
		"recursive":                           false,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", DROPBOX_API_URL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Dropbox API returns folders in an "entries" field
	type DropboxAPIResponse struct {
		Entries []DropboxFolder `json:"entries"`
	}

	var apiResponse DropboxAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	// Convert to DropboxFolders type
	folders := DropboxFolders(apiResponse.Entries)

	ctx.Logger.Info("Successfully retrieved dropbox folders", "count", len(folders))
	return folders, nil
}
