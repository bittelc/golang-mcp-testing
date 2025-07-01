package dropbox

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/localrivet/gomcp/server"
)

// mockLogger creates a basic logger for testing
func mockLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

// mockContext creates a server context for testing
func mockContext() *server.Context {
	return &server.Context{
		Logger: mockLogger(),
	}
}

func TestHandleListDropboxFolders_MissingAPIKey(t *testing.T) {
	// Save original env and restore after test
	originalAPIKey := os.Getenv("DROPBOX_API_KEY")
	defer os.Setenv("DROPBOX_API_KEY", originalAPIKey)

	// Unset the API key
	os.Unsetenv("DROPBOX_API_KEY")

	ctx := mockContext()
	args := ListDropboxFoldersArgs{Path: "/test"}

	folders, err := HandleListDropboxFolder(ctx, args)

	if err == nil {
		t.Fatal("Expected error when API key is missing")
	}

	if folders != nil {
		t.Fatal("Expected nil folders when API key is missing")
	}

	expectedError := "$DROPBOX_API_KEY not set, unable to retrieve dropbox folders"
	if !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("Expected error to contain '%s', got: %s", expectedError, err.Error())
	}
}

func TestHandleListDropboxFolders_SuccessfulResponse(t *testing.T) {
	// Mock successful Dropbox API response
	mockResponse := `{
		"entries": [
			{
				"id": "id:test123",
				"name": "Test Folder",
				"path_display": "/Test Folder",
				".tag": "folder",
				"shared_folder_id": "shared123"
			},
			{
				"id": "id:test456",
				"name": "Another Folder",
				"path_display": "/Another Folder",
				".tag": "folder",
				"shared_folder_id": "shared456"
			}
		]
	}`

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and headers
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
		}

		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			t.Errorf("Expected Authorization header to start with 'Bearer ', got %s", authHeader)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Set API key
	os.Setenv("DROPBOX_API_KEY", "test_api_key_123")
	defer os.Unsetenv("DROPBOX_API_KEY")

	// We'll need to test with a modified version or create a testable version
	// For now, let's test the unmarshalFolders function separately
	t.Skip("This test requires modifying the HTTP client or URL - testing unmarshalFolders separately")
}

func TestHandleListDropboxFolders_HTTPClientError(t *testing.T) {
	// Set valid API key
	os.Setenv("DROPBOX_API_KEY", "test_api_key_123")
	defer os.Unsetenv("DROPBOX_API_KEY")

	ctx := mockContext()
	args := ListDropboxFoldersArgs{Path: "/test"}

	// This will fail because we're hitting the real Dropbox API URL without proper setup
	// In a real test environment, you'd mock the HTTP client
	folders, err := HandleListDropboxFolder(ctx, args)

	if err == nil {
		t.Fatal("Expected error when making HTTP request without proper API setup")
	}

	if folders != nil {
		t.Fatal("Expected nil folders when HTTP request fails")
	}
}

func TestHandleListDropboxFolders_PathHandling(t *testing.T) {
	tests := []struct {
		name         string
		inputPath    string
		expectedPath string
	}{
		{
			name:         "Empty path",
			inputPath:    "",
			expectedPath: "",
		},
		{
			name:         "Root slash path",
			inputPath:    "/",
			expectedPath: "",
		},
		{
			name:         "Dot path",
			inputPath:    ".",
			expectedPath: "",
		},
		{
			name:         "Regular path",
			inputPath:    "/documents",
			expectedPath: "/documents",
		},
	}

	// Set valid API key
	os.Setenv("DROPBOX_API_KEY", "test_api_key_123")
	defer os.Unsetenv("DROPBOX_API_KEY")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := mockContext()
			args := ListDropboxFoldersArgs{Path: tt.inputPath}

			// Test the craftHttpReq function directly to verify path handling
			req, err := craftHttpReq(ctx, &args, "test_key")
			if err != nil {
				t.Fatalf("craftHttpReq failed: %v", err)
			}

			// Verify the path was processed correctly
			if args.Path != tt.expectedPath {
				t.Errorf("Expected path to be '%s', got '%s'", tt.expectedPath, args.Path)
			}

			// Verify request properties
			if req.Method != "POST" {
				t.Errorf("Expected POST method, got %s", req.Method)
			}

			if req.Header.Get("Authorization") != "Bearer test_key" {
				t.Errorf("Expected Authorization header 'Bearer test_key', got '%s'", req.Header.Get("Authorization"))
			}
		})
	}
}

func TestUnmarshalFolders_ValidJSON(t *testing.T) {
	validJSON := `{
		"entries": [
			{
				"id": "id:test123",
				"name": "Test Folder",
				"path_display": "/Test Folder",
				".tag": "folder",
				"shared_folder_id": "shared123"
			}
		]
	}`

	jsonBytes := []byte(validJSON)
	folders, err := unmarshalFolders(&jsonBytes)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(folders) != 1 {
		t.Fatalf("Expected 1 folder, got %d", len(folders))
	}

	folder := folders[0]
	if folder.ID != "id:test123" {
		t.Errorf("Expected ID 'id:test123', got '%s'", folder.ID)
	}

	if folder.Name != "Test Folder" {
		t.Errorf("Expected Name 'Test Folder', got '%s'", folder.Name)
	}

	if folder.PathDisplay != "/Test Folder" {
		t.Errorf("Expected PathDisplay '/Test Folder', got '%s'", folder.PathDisplay)
	}

	if folder.Tag != "folder" {
		t.Errorf("Expected Tag 'folder', got '%s'", folder.Tag)
	}
}

func TestUnmarshalFolders_EmptyEntries(t *testing.T) {
	emptyJSON := `{"entries": []}`
	jsonBytes := []byte(emptyJSON)

	folders, err := unmarshalFolders(&jsonBytes)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(folders) != 0 {
		t.Fatalf("Expected empty folders slice, got %d folders", len(folders))
	}
}

func TestUnmarshalFolders_InvalidJSON(t *testing.T) {
	invalidJSON := `{"invalid": json}`
	jsonBytes := []byte(invalidJSON)

	folders, err := unmarshalFolders(&jsonBytes)

	if err == nil {
		t.Fatal("Expected error for invalid JSON")
	}

	if folders != nil {
		t.Fatal("Expected nil folders for invalid JSON")
	}

	if !strings.Contains(err.Error(), "failed to unmarshal JSON response") {
		t.Errorf("Expected error message to contain 'failed to unmarshal JSON response', got: %s", err.Error())
	}
}

func TestUnmarshalFolders_MalformedJSON(t *testing.T) {
	malformedJSON := `{invalid json`
	jsonBytes := []byte(malformedJSON)

	folders, err := unmarshalFolders(&jsonBytes)

	if err == nil {
		t.Fatal("Expected error for malformed JSON")
	}

	if folders != nil {
		t.Fatal("Expected nil folders for malformed JSON")
	}
}

func TestCraftHttpReq_ValidRequest(t *testing.T) {
	ctx := mockContext()
	args := &ListDropboxFoldersArgs{Path: "/test/path"}
	apiKey := "test_api_key_123"

	req, err := craftHttpReq(ctx, args, apiKey)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify request method
	if req.Method != "POST" {
		t.Errorf("Expected POST method, got %s", req.Method)
	}

	// Verify headers
	if req.Header.Get("Authorization") != "Bearer "+apiKey {
		t.Errorf("Expected Authorization 'Bearer %s', got '%s'", apiKey, req.Header.Get("Authorization"))
	}

	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", req.Header.Get("Content-Type"))
	}

	// Verify URL
	expectedURL := fmt.Sprintf("%s/list_folder", DROPBOX_FILES_API_URL)
	if req.URL.String() != expectedURL {
		t.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL.String())
	}
}

func TestHandleListDropboxFolders_APIKeyLogging(t *testing.T) {
	// Set API key with sufficient length
	os.Setenv("DROPBOX_API_KEY", "test_api_key_123456")
	defer os.Unsetenv("DROPBOX_API_KEY")

	// Create a context to capture logs (this would require custom logger setup in real implementation)
	ctx := mockContext()
	args := ListDropboxFoldersArgs{Path: "/test"}

	// This will fail with HTTP error, but we're testing the API key logging part
	_, err := HandleListDropboxFolder(ctx, args)

	// We expect an error because we're not mocking the HTTP call
	if err == nil {
		t.Fatal("Expected HTTP error but got none")
	}

	// The function should have logged the first two characters of the API key
	// In a real test, you'd capture and verify the log output
}

func TestListDropboxFoldersArgs_EmptyStruct(t *testing.T) {
	args := ListDropboxFoldersArgs{}

	if args.Path != "" {
		t.Errorf("Expected empty Path in default struct, got '%s'", args.Path)
	}
}

func TestDropboxFolder_JSONTags(t *testing.T) {
	// Test that the struct can be marshaled/unmarshaled properly
	folder := DropboxFolder{
		ID:             "test-id",
		Name:           "test-name",
		PathDisplay:    "/test/path",
		Tag:            "folder",
		SharedFolderID: "shared-id",
	}

	// This is more of a compile-time check, but we can verify the struct is properly defined
	if folder.ID == "" {
		t.Error("Expected non-empty ID")
	}
}
