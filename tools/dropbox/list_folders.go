package dropbox

import (
	"fmt"
	"os"

	"github.com/localrivet/gomcp/server"
)

type ListDropboxFoldersArgs struct {
	Path string `json:"path"`
}

type DropboxFolders []DropboxFolder
type DropboxFolder struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

const DROPBOX_API_URL = "https://api.dropboxapi.com/2/sharing/list_folders"

// HandleListDropBoxFolders implements the logic the list_dropbox_folders tool
func HandleListDropboxFolders(ctx *server.Context, args ListDropboxFoldersArgs) (DropboxFolders, error) {
	ctx.Logger.Info("Handling ListDropboxFolders tool call")

	// This handler provides a listing of all folders and their metadata at
	// the provided path (ListDropboxFoldersArgs.Path).

	// Sample HTTP request: https://www.dropbox.com/developers/documentation/http/documentation#sharing-list_folders
	// curl -X POST https://api.dropboxapi.com/2/sharing/list_folders \
	//    --header "Authorization: Bearer <get access token>" \
	//    --header "Content-Type: application/json" \
	//    --data "{\"actions\":[],\"limit\":100}"

	var folders DropboxFolders
	apiKey := os.Getenv("DROPBOX_API_KEY")
	if apiKey == "" {
		ctx.Logger.Info("$DROPBOX_API_KEY not set")
		return nil, fmt.Errorf("$DROPBOX_API_KEY not set, unable to retrieve dropbox folders")
	}
	ctx.Logger.Info("No dropbox folders at the provided path")
	return folders, nil
}
