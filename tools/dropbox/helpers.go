package dropbox

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const DROPBOX_API_URL = "https://api.dropboxapi.com/2/files/list_folder"

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
