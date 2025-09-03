package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/backendArchitect/go-api-gen/pkg/generator"
	"github.com/backendArchitect/go-api-gen/pkg/logger"
)

var (
	inputFile    string
	outputDir    string
	packageName  string
	clientName   string
	logLevel     string
	debug        bool
	jsonLog      bool
)

var rootCmd = &cobra.Command{
	Use:   "go-api-gen",
	Short: "Generate Go client libraries from OpenAPI/Swagger specifications",
	Long: `go-api-gen is a tool that takes an OpenAPI/Swagger specification and 
automatically generates a complete Go client library for that API. This saves 
developers time and effort when integrating with a new service.

SUPPORTED FORMATS:
- JSON (.json files) - Standard Swagger/OpenAPI JSON format
- YAML (.yaml/.yml files) - YAML format OpenAPI specifications

LOGGING:
The tool provides comprehensive logging for debugging and monitoring:
- Use --log-level to control verbosity (debug, info, warn, error)
- Use --debug for detailed debug output  
- Use --json-log for structured JSON logging
- Components are logged separately (cli, generator, parser, templates)

EXAMPLES:
  # Generate from YAML specification
  go-api-gen --input petstore.yaml --output ./generated
  
  # Generate from JSON specification  
  go-api-gen --input swagger.json --output ./generated
  
  # Generate with debug logging for troubleshooting
  go-api-gen --input petstore.json --output ./generated --debug
  
  # Generate with JSON logging for monitoring systems
  go-api-gen --input swagger.json --output ./generated --json-log`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateClient()
	},
}

func init() {
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Path to OpenAPI/Swagger specification file (JSON or YAML format) (required)")
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", "./generated", "Output directory for generated client code")
	rootCmd.Flags().StringVarP(&packageName, "package", "p", "client", "Go package name for generated code")
	rootCmd.Flags().StringVarP(&clientName, "client-name", "c", "APIClient", "Name for the generated client struct")
	rootCmd.Flags().StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error)")
	rootCmd.Flags().BoolVar(&debug, "debug", false, "Enable debug logging (equivalent to --log-level=debug)")
	rootCmd.Flags().BoolVar(&jsonLog, "json-log", false, "Output logs in JSON format")
	
	rootCmd.MarkFlagRequired("input")
}

func Execute() error {
	return rootCmd.Execute()
}

func generateClient() error {
	// Configure logging
	loggerConfig := logger.Config{
		Level:  logger.ParseLevel(logLevel),
		Format: "text",
	}
	
	// Override with debug flag if set
	if debug {
		loggerConfig.Level = logger.DebugLevel
	}
	
	// Use JSON format if requested
	if jsonLog {
		loggerConfig.Format = "json"
	}
	
	log := logger.New(loggerConfig).WithComponent("cli")
	
	log.InfoContext("Starting code generation", 
		"input_file", inputFile,
		"output_dir", outputDir,
		"package_name", packageName,
		"client_name", clientName)
	
	log.DebugContext("Configuration details",
		"log_level", logLevel,
		"debug", debug,
		"json_log", jsonLog)
	
	// Validate input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		log.ErrorContext("Input file does not exist", "file", inputFile)
		return fmt.Errorf("input file does not exist: %s", inputFile)
	}
	
	log.DebugContext("Input file validated", "file", inputFile)

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.ErrorContext("Failed to create output directory", "dir", outputDir, "error", err)
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	log.DebugContext("Output directory prepared", "dir", outputDir)

	// Initialize generator
	gen := generator.New(generator.Config{
		InputFile:   inputFile,
		OutputDir:   outputDir,
		PackageName: packageName,
		ClientName:  clientName,
		Logger:      log,
	})

	// Generate the client
	if err := gen.Generate(); err != nil {
		log.ErrorContext("Code generation failed", "error", err)
		return fmt.Errorf("failed to generate client: %w", err)
	}

	log.InfoContext("Code generation completed successfully", "output_dir", outputDir)
	fmt.Printf("Successfully generated Go client in %s\n", outputDir)
	return nil
}