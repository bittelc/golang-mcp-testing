package dropbox

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/localrivet/gomcp/server"
)

const DROPBOX_FILES_API_URL = "https://api.dropboxapi.com/2/files"

func handleFailedHttpReq(resp *http.Response) error {
	// Read the response body to get more details about the error
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("API request FAILED with status: %d, failed to read error response: %w", resp.StatusCode, err)
	}

	// Try to parse error response as JSON for structured error info
	var errorResponse map[string]any
	if json.Unmarshal(body, &errorResponse) == nil {
		return fmt.Errorf("API request FAILURE with status: %d, error: %v", resp.StatusCode, errorResponse)
	}

	// If JSON parsing fails, include raw response body
	return fmt.Errorf("API request FAILING with status: %d, response: %s", resp.StatusCode, string(body))
}

func HelperCallHandlerDirectly(logger *slog.Logger, functionName string, args interface{}) error {
	serverContext := &server.Context{
		Logger: logger,
	}

	switch functionName {
	case "HandleListDropboxFolders":
		folderArgs, ok := args.(ListDropboxFoldersArgs)
		if !ok {
			return fmt.Errorf("invalid arguments for HandleListDropboxFolders, expected ListDropboxFoldersArgs")
		}
		folders, err := HandleListDropboxFolders(serverContext, folderArgs)
		if err != nil {
			log.Fatalf("Failed to list dropbox folders, %v", err)
			return err
		}
		spew.Dump(folders)

	case "HandleFilesDownload":
		downloadArgs, ok := args.(FilesDownloadArgs)
		if !ok {
			return fmt.Errorf("invalid arguments for HandleFilesDownload, expected FilesDownloadArgs")
		}
		metadata, err := HandleFilesDownload(serverContext, downloadArgs)
		if err != nil {
			log.Fatalf("Failed to download file, %v", err)
			return err
		}
		spew.Dump(metadata)

	default:
		return fmt.Errorf("unknown function name: %s", functionName)
	}

	return nil
}
