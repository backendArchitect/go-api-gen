package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/backendArchitect/go-api-gen/pkg/logger"
	"github.com/backendArchitect/go-api-gen/pkg/parser"
	"github.com/backendArchitect/go-api-gen/pkg/templates"
)

// Config holds the configuration for the generator
type Config struct {
	InputFile   string
	OutputDir   string
	PackageName string
	ClientName  string
	Logger      *logger.Logger
}

// Generator handles the code generation process
type Generator struct {
	config Config
	parser *parser.Parser
	logger *logger.Logger
}

// New creates a new Generator instance
func New(config Config) *Generator {
	// Create a default logger if none provided
	if config.Logger == nil {
		config.Logger = logger.New(logger.Config{Level: logger.InfoLevel})
	}
	
	return &Generator{
		config: config,
		parser: parser.New(config.Logger),
		logger: config.Logger.WithComponent("generator"),
	}
}

// Generate parses the OpenAPI spec and generates Go client code
func (g *Generator) Generate() error {
	g.logger.DebugContext("Starting code generation process")
	
	// Create output directory if it doesn't exist
	g.logger.DebugContext("Creating output directory", "dir", g.config.OutputDir)
	if err := os.MkdirAll(g.config.OutputDir, 0755); err != nil {
		g.logger.ErrorContext("Failed to create output directory", "dir", g.config.OutputDir, "error", err)
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Parse the OpenAPI specification
	g.logger.InfoContext("Parsing OpenAPI specification", "file", g.config.InputFile)
	spec, err := g.parser.Parse(g.config.InputFile)
	if err != nil {
		g.logger.ErrorContext("Failed to parse OpenAPI spec", "file", g.config.InputFile, "error", err)
		return fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}
	g.logger.InfoContext("Successfully parsed OpenAPI specification")

	// Generate client code
	g.logger.InfoContext("Generating client code")
	clientCode, err := templates.GenerateClient(spec, g.config.PackageName, g.config.ClientName, g.logger)
	if err != nil {
		g.logger.ErrorContext("Failed to generate client code", "error", err)
		return fmt.Errorf("failed to generate client code: %w", err)
	}
	g.logger.DebugContext("Client code generation completed")

	// Generate model code
	g.logger.InfoContext("Generating model code")
	modelCode, err := templates.GenerateModels(spec, g.config.PackageName, g.logger)
	if err != nil {
		g.logger.ErrorContext("Failed to generate model code", "error", err)
		return fmt.Errorf("failed to generate model code: %w", err)
	}
	g.logger.DebugContext("Model code generation completed")

	// Write files
	g.logger.InfoContext("Writing generated files")
	if err := g.writeFile("client.go", clientCode); err != nil {
		return err
	}

	if err := g.writeFile("models.go", modelCode); err != nil {
		return err
	}

	g.logger.InfoContext("Code generation completed successfully", 
		"files_written", []string{"client.go", "models.go"},
		"output_dir", g.config.OutputDir)
	return nil
}

// writeFile writes content to a file in the output directory
func (g *Generator) writeFile(filename, content string) error {
	path := filepath.Join(g.config.OutputDir, filename)
	g.logger.DebugContext("Writing file", "file", path, "size_bytes", len(content))
	
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		g.logger.ErrorContext("Failed to write file", "file", path, "error", err)
		return err
	}
	
	g.logger.InfoContext("File written successfully", "file", filename, "path", path)
	return nil
}