package parser

import (
	"context"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/backendArchitect/go-api-gen/pkg/logger"
)

// Parser handles parsing of OpenAPI specifications
type Parser struct {
	logger *logger.Logger
}

// New creates a new Parser instance
func New(log *logger.Logger) *Parser {
	// Create a default logger if none provided
	if log == nil {
		log = logger.New(logger.Config{Level: logger.InfoLevel})
	}
	
	return &Parser{
		logger: log.WithComponent("parser"),
	}
}

// Parse parses an OpenAPI specification file and returns the parsed document
func (p *Parser) Parse(filename string) (*openapi3.T, error) {
	p.logger.DebugContext("Loading OpenAPI specification", "file", filename)
	
	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	
	// Load the OpenAPI document
	doc, err := loader.LoadFromFile(filename)
	if err != nil {
		p.logger.ErrorContext("Failed to load OpenAPI document", "file", filename, "error", err)
		return nil, fmt.Errorf("failed to load OpenAPI document: %w", err)
	}
	
	p.logger.DebugContext("OpenAPI document loaded successfully", "file", filename)

	// Validate the document
	p.logger.DebugContext("Validating OpenAPI document")
	if err := doc.Validate(ctx); err != nil {
		p.logger.ErrorContext("Invalid OpenAPI document", "file", filename, "error", err)
		return nil, fmt.Errorf("invalid OpenAPI document: %w", err)
	}
	
	// Count operations for logging
	operationCount := 0
	for _, pathItem := range doc.Paths.Map() {
		if pathItem.Get != nil { operationCount++ }
		if pathItem.Post != nil { operationCount++ }
		if pathItem.Put != nil { operationCount++ }
		if pathItem.Delete != nil { operationCount++ }
		if pathItem.Patch != nil { operationCount++ }
		if pathItem.Options != nil { operationCount++ }
		if pathItem.Head != nil { operationCount++ }
		if pathItem.Trace != nil { operationCount++ }
	}
	
	p.logger.InfoContext("OpenAPI document parsed and validated successfully", 
		"file", filename,
		"title", doc.Info.Title,
		"version", doc.Info.Version,
		"operation_count", operationCount,
		"server_count", len(doc.Servers))

	return doc, nil
}