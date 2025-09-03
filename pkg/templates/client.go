package templates

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/backendArchitect/go-api-gen/pkg/logger"
)

// ClientData holds the data needed for client template generation
type ClientData struct {
	PackageName string
	ClientName  string
	BaseURL     string
	Operations  []Operation
}

// Operation represents an API operation
type Operation struct {
	Name        string
	Method      string
	Path        string
	Summary     string
	RequestType string
	ResponseType string
	Parameters  []Parameter
}

// Parameter represents an operation parameter
type Parameter struct {
	Name     string
	Type     string
	Required bool
	In       string // query, path, header
}

// GenerateClient generates the main client code
func GenerateClient(spec *openapi3.T, packageName, clientName string, log *logger.Logger) (string, error) {
	// Create a default logger if none provided
	if log == nil {
		log = logger.New(logger.Config{Level: logger.InfoLevel})
	}
	clientLogger := log.WithComponent("client-template")
	
	clientLogger.DebugContext("Starting client code generation",
		"package_name", packageName,
		"client_name", clientName)
	
	data := ClientData{
		PackageName: packageName,
		ClientName:  clientName,
		BaseURL:     getBaseURL(spec, clientLogger),
		Operations:  extractOperations(spec, clientLogger),
	}
	
	clientLogger.DebugContext("Template data prepared", 
		"operation_count", len(data.Operations),
		"base_url", data.BaseURL)

	tmpl := template.Must(template.New("client").Parse(clientTemplate))
	
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		clientLogger.ErrorContext("Failed to execute client template", "error", err)
		return "", fmt.Errorf("failed to execute client template: %w", err)
	}
	
	clientLogger.DebugContext("Client template executed successfully", "size_bytes", buf.Len())
	return buf.String(), nil
}

// getBaseURL extracts the base URL from the OpenAPI spec
func getBaseURL(spec *openapi3.T, log *logger.Logger) string {
	if len(spec.Servers) > 0 {
		baseURL := spec.Servers[0].URL
		log.DebugContext("Using base URL from spec", "url", baseURL)
		return baseURL
	}
	
	defaultURL := "https://api.example.com"
	log.WarnContext("No servers defined in spec, using default", "default_url", defaultURL)
	return defaultURL
}

// extractOperations extracts all operations from the OpenAPI spec
func extractOperations(spec *openapi3.T, log *logger.Logger) []Operation {
	log.DebugContext("Starting operation extraction")
	var operations []Operation

	// Sort paths for consistent output
	paths := make([]string, 0, len(spec.Paths.Map()))
	for path := range spec.Paths.Map() {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	
	log.DebugContext("Processing paths", "path_count", len(paths))

	for _, path := range paths {
		pathItem := spec.Paths.Map()[path]
		
		// Check each HTTP method
		methods := map[string]*openapi3.Operation{
			"GET":    pathItem.Get,
			"POST":   pathItem.Post,
			"PUT":    pathItem.Put,
			"DELETE": pathItem.Delete,
			"PATCH":  pathItem.Patch,
		}

		for method, op := range methods {
			if op != nil {
				operation := Operation{
					Name:        generateOperationName(method, path, op),
					Method:      method,
					Path:        path,
					Summary:     op.Summary,
					Parameters:  extractParameters(op, log),
				}
				
				log.DebugContext("Extracted operation", 
					"name", operation.Name,
					"method", method,
					"path", path,
					"parameter_count", len(operation.Parameters))
				
				operations = append(operations, operation)
			}
		}
	}
	
	log.InfoContext("Operation extraction completed", "total_operations", len(operations))
	return operations
}

// generateOperationName creates a Go method name from the operation
func generateOperationName(method, path string, op *openapi3.Operation) string {
	if op.OperationID != "" {
		return toCamelCase(op.OperationID)
	}
	
	// Generate name from method and path
	parts := strings.Split(path, "/")
	var nameParts []string
	nameParts = append(nameParts, strings.ToLower(method))
	
	for _, part := range parts {
		if part != "" && !strings.HasPrefix(part, "{") {
			nameParts = append(nameParts, part)
		}
	}
	
	return toCamelCase(strings.Join(nameParts, "_"))
}

// extractParameters extracts parameters from an operation
func extractParameters(op *openapi3.Operation, log *logger.Logger) []Parameter {
	var params []Parameter
	
	for _, paramRef := range op.Parameters {
		if paramRef.Value != nil {
			param := Parameter{
				Name:     paramRef.Value.Name,
				Type:     getParameterType(paramRef.Value),
				Required: paramRef.Value.Required,
				In:       paramRef.Value.In,
			}
			params = append(params, param)
			
			log.DebugContext("Extracted parameter",
				"name", param.Name,
				"type", param.Type,
				"required", param.Required,
				"in", param.In)
		}
	}
	
	return params
}

// getParameterType determines the Go type for a parameter
func getParameterType(param *openapi3.Parameter) string {
	if param.Schema != nil && param.Schema.Value != nil {
		return openAPITypeToGo(param.Schema.Value.Type.Slice())
	}
	return "string"
}

// openAPITypeToGo converts OpenAPI types to Go types
func openAPITypeToGo(types []string) string {
	if len(types) == 0 {
		return "string"
	}
	
	switch types[0] {
	case "integer":
		return "int"
	case "number":
		return "float64"
	case "boolean":
		return "bool"
	case "array":
		return "[]interface{}"
	case "object":
		return "map[string]interface{}"
	default:
		return "string"
	}
}