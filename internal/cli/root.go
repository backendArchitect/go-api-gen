package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/backendArchitect/go-api-gen/pkg/generator"
)

var (
	inputFile    string
	outputDir    string
	packageName  string
	clientName   string
)

var rootCmd = &cobra.Command{
	Use:   "go-api-gen",
	Short: "Generate Go client libraries from OpenAPI/Swagger specifications",
	Long: `go-api-gen is a tool that takes an OpenAPI/Swagger specification and 
automatically generates a complete Go client library for that API. This saves 
developers time and effort when integrating with a new service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateClient()
	},
}

func init() {
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Path to OpenAPI/Swagger specification file (required)")
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", "./generated", "Output directory for generated client code")
	rootCmd.Flags().StringVarP(&packageName, "package", "p", "client", "Go package name for generated code")
	rootCmd.Flags().StringVarP(&clientName, "client-name", "c", "APIClient", "Name for the generated client struct")
	
	rootCmd.MarkFlagRequired("input")
}

func Execute() error {
	return rootCmd.Execute()
}

func generateClient() error {
	// Validate input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", inputFile)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Initialize generator
	gen := generator.New(generator.Config{
		InputFile:   inputFile,
		OutputDir:   outputDir,
		PackageName: packageName,
		ClientName:  clientName,
	})

	// Generate the client
	if err := gen.Generate(); err != nil {
		return fmt.Errorf("failed to generate client: %w", err)
	}

	fmt.Printf("Successfully generated Go client in %s\n", outputDir)
	return nil
}