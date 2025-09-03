package templates

import (
	"bytes"
	"fmt"
	"sort"
	"text/template"

	"github.com/getkin/kin-openapi/openapi3"
)

// ModelData holds the data needed for model template generation
type ModelData struct {
	PackageName string
	Models      []Model
	TypeAliases []TypeAlias
	NeedsTime   bool
}

// Model represents a data model
type Model struct {
	Name        string
	Description string
	Fields      []Field
}

// TypeAlias represents a type alias
type TypeAlias struct {
	Name        string
	Type        string
	Description string
}

// Field represents a struct field
type Field struct {
	Name        string
	Type        string
	JSONTag     string
	Description string
	Required    bool
}

// GenerateModels generates the model structs code
func GenerateModels(spec *openapi3.T, packageName string) (string, error) {
	models, typeAliases := extractModelsAndAliases(spec)
	
	// Check if we need time import
	needsTime := false
	for _, model := range models {
		for _, field := range model.Fields {
			if field.Type == "time.Time" {
				needsTime = true
				break
			}
		}
		if needsTime {
			break
		}
	}
	
	data := ModelData{
		PackageName: packageName,
		Models:      models,
		TypeAliases: typeAliases,
		NeedsTime:   needsTime,
	}

	tmpl := template.Must(template.New("models").Parse(modelsTemplate))
	
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute models template: %w", err)
	}

	return buf.String(), nil
}

// extractModelsAndAliases extracts all data models and type aliases from the OpenAPI spec
func extractModelsAndAliases(spec *openapi3.T) ([]Model, []TypeAlias) {
	var models []Model
	var typeAliases []TypeAlias

	if spec.Components == nil || spec.Components.Schemas == nil {
		return models, typeAliases
	}

	// Sort schema names for consistent output
	names := make([]string, 0, len(spec.Components.Schemas))
	for name := range spec.Components.Schemas {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		schemaRef := spec.Components.Schemas[name]
		if schemaRef.Value != nil {
			types := schemaRef.Value.Type.Slice()
			
			// Handle array types as type aliases
			if len(types) > 0 && types[0] == "array" {
				alias := TypeAlias{
					Name:        toCamelCase(name),
					Type:        schemaToGoType(schemaRef.Value),
					Description: schemaRef.Value.Description,
				}
				typeAliases = append(typeAliases, alias)
			} else {
				// Regular object types as structs
				model := Model{
					Name:        toCamelCase(name),
					Description: schemaRef.Value.Description,
					Fields:      extractFields(schemaRef.Value),
				}
				models = append(models, model)
			}
		}
	}

	return models, typeAliases
}

// extractFields extracts fields from a schema
func extractFields(schema *openapi3.Schema) []Field {
	var fields []Field

	if schema.Properties == nil {
		return fields
	}

	// Sort property names for consistent output
	names := make([]string, 0, len(schema.Properties))
	for name := range schema.Properties {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		propRef := schema.Properties[name]
		if propRef.Value != nil {
			field := Field{
				Name:        toCamelCase(name),
				Type:        schemaToGoType(propRef.Value),
				JSONTag:     name,
				Description: propRef.Value.Description,
				Required:    isRequired(name, schema.Required),
			}
			fields = append(fields, field)
		}
	}

	return fields
}

// schemaToGoType converts an OpenAPI schema to a Go type
func schemaToGoType(schema *openapi3.Schema) string {
	types := schema.Type.Slice()
	if len(types) == 0 {
		return "interface{}"
	}
	
	switch types[0] {
	case "string":
		if schema.Format == "date-time" {
			return "time.Time"
		}
		return "string"
	case "integer":
		if schema.Format == "int64" {
			return "int64"
		}
		return "int"
	case "number":
		if schema.Format == "float" {
			return "float32"
		}
		return "float64"
	case "boolean":
		return "bool"
	case "array":
		if schema.Items != nil && schema.Items.Ref != "" {
			// Handle $ref in array items first
			refName := extractRefName(schema.Items.Ref)
			return "[]" + toCamelCase(refName)
		} else if schema.Items != nil && schema.Items.Value != nil {
			itemType := schemaToGoType(schema.Items.Value)
			return "[]" + itemType
		}
		return "[]interface{}"
	case "object":
		if len(schema.Properties) == 0 {
			return "map[string]interface{}"
		}
		// For complex objects, we'd need to generate nested structs
		// For now, fallback to map
		return "map[string]interface{}"
	default:
		return "interface{}"
	}
}

// isRequired checks if a field is required
func isRequired(fieldName string, required []string) bool {
	for _, req := range required {
		if req == fieldName {
			return true
		}
	}
	return false
}