package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/backendArchitect/go-api-gen/pkg/parser"
	"github.com/backendArchitect/go-api-gen/pkg/templates"
)

// Config holds the configuration for the generator
type Config struct {
	InputFile   string
	OutputDir   string
	PackageName string
	ClientName  string
}

// Generator handles the code generation process
type Generator struct {
	config Config
	parser *parser.Parser
}

// New creates a new Generator instance
func New(config Config) *Generator {
	return &Generator{
		config: config,
		parser: parser.New(),
	}
}

// Generate parses the OpenAPI spec and generates Go client code
func (g *Generator) Generate() error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(g.config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Parse the OpenAPI specification
	spec, err := g.parser.Parse(g.config.InputFile)
	if err != nil {
		return fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	// Generate client code
	clientCode, err := templates.GenerateClient(spec, g.config.PackageName, g.config.ClientName)
	if err != nil {
		return fmt.Errorf("failed to generate client code: %w", err)
	}

	// Generate model code
	modelCode, err := templates.GenerateModels(spec, g.config.PackageName)
	if err != nil {
		return fmt.Errorf("failed to generate model code: %w", err)
	}

	// Write files
	if err := g.writeFile("client.go", clientCode); err != nil {
		return err
	}

	if err := g.writeFile("models.go", modelCode); err != nil {
		return err
	}

	return nil
}

// writeFile writes content to a file in the output directory
func (g *Generator) writeFile(filename, content string) error {
	path := filepath.Join(g.config.OutputDir, filename)
	return os.WriteFile(path, []byte(content), 0644)
}