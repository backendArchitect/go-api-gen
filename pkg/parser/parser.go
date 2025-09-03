package parser

import (
	"context"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// Parser handles parsing of OpenAPI specifications
type Parser struct{}

// New creates a new Parser instance
func New() *Parser {
	return &Parser{}
}

// Parse parses an OpenAPI specification file and returns the parsed document
func (p *Parser) Parse(filename string) (*openapi3.T, error) {
	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	
	// Load the OpenAPI document
	doc, err := loader.LoadFromFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load OpenAPI document: %w", err)
	}

	// Validate the document
	if err := doc.Validate(ctx); err != nil {
		return nil, fmt.Errorf("invalid OpenAPI document: %w", err)
	}

	return doc, nil
}