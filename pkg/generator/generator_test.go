package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/backendArchitect/go-api-gen/pkg/logger"
)

func TestGenerateFromPetstore(t *testing.T) {
	// Create temporary directory for test output
	tmpDir := filepath.Join(os.TempDir(), "go-api-gen-test")
	defer os.RemoveAll(tmpDir)

	// Create a test logger
	testLogger := logger.New(logger.Config{Level: logger.DebugLevel})

	config := Config{
		InputFile:   "../../examples/petstore.yaml",
		OutputDir:   tmpDir,
		PackageName: "client",
		ClientName:  "APIClient",
		Logger:      testLogger,
	}

	gen := New(config)
	
	if err := gen.Generate(); err != nil {
		t.Fatalf("Failed to generate client: %v", err)
	}

	// Check that files were created
	clientFile := filepath.Join(tmpDir, "client.go")
	modelsFile := filepath.Join(tmpDir, "models.go")

	if _, err := os.Stat(clientFile); os.IsNotExist(err) {
		t.Errorf("client.go was not generated")
	}

	if _, err := os.Stat(modelsFile); os.IsNotExist(err) {
		t.Errorf("models.go was not generated")
	}
}