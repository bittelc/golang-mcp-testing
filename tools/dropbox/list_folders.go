package dropbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
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

/*
 * {
   "entries": [
     {
       ".tag": "folder",
       "name": "Family Room",
       "path_lower": "/family room",
       "path_display": "/Family Room",
       "id": "id:cMZquz1z9UkAAAAAAAEusw",
       "shared_folder_id": "2110133921",
       "sharing_info": {
         "read_only": false,
         "shared_folder_id": "2110133921",
         "traverse_only": false,
         "no_access": false
       }
     },
     {
       ".tag": "folder",
       "name": "Camera Uploads",
       "path_lower": "/camera uploads",
       "path_display": "/Camera Uploads",
       "id": "id:cMZquz1z9UkAAAAAAAE1aw",
       "shared_folder_id": "3750858913",
       "sharing_info": {
         "read_only": false,
         "shared_folder_id": "3750858913",
         "traverse_only": false,
         "no_access": false
       }
     },
     {
       ".tag": "folder",
       "name": "Cole Personal",
       "path_lower": "/cole personal",
       "path_display": "/Cole Personal",
       "id": "id:cMZquz1z9UkAAAAAAAE5SA",
       "shared_folder_id": "3871520337",
       "sharing_info": {
         "read_only": false,
         "shared_folder_id": "3871520337",
         "traverse_only": false,
         "no_access": false
       }
     },
     {
       ".tag": "folder",
       "name": "Apps",
       "path_lower": "/apps",
       "path_display": "/Apps",
       "id": "id:cMZquz1z9UkAAAAAAAJTrw"
     }
   ],
*/

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
		"include_deleted":                     true,
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
	spew.Dump(resp)
	var folders DropboxFolders
	ctx.Logger.Info("No dropbox folders at the provided path")
	return folders, nil
}
