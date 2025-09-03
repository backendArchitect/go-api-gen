# go-api-gen

A Go tool that automatically generates complete Go client libraries from OpenAPI/Swagger specifications. This saves developers time and effort when integrating with new APIs by automatically creating all the necessary HTTP client code, request/response models, and error handling.

## Features

- **OpenAPI 3.0 Support**: Parses OpenAPI 3.0 specifications (JSON/YAML)
- **Complete Client Generation**: Generates HTTP client with methods for all API endpoints
- **Model Generation**: Creates Go structs for all data models with proper JSON tags
- **Type Safety**: Converts OpenAPI types to appropriate Go types
- **Error Handling**: Built-in HTTP error handling and response parsing
- **Customizable**: Configurable package names, client names, and output directories

## Installation

```bash
go install github.com/backendArchitect/go-api-gen/cmd/go-api-gen@latest
```

Or build from source:

```bash
git clone https://github.com/backendArchitect/go-api-gen.git
cd go-api-gen
go build ./cmd/go-api-gen
```

## Usage

### Basic Usage

```bash
go-api-gen --input petstore.yaml --output ./generated
```

### Advanced Usage

```bash
go-api-gen \
  --input examples/petstore.yaml \
  --output ./clients/petstore \
  --package petstore \
  --client-name PetStoreClient
```

### Command Line Options

- `--input`, `-i`: Path to OpenAPI/Swagger specification file (required)
- `--output`, `-o`: Output directory for generated client code (default: `./generated`)
- `--package`, `-p`: Go package name for generated code (default: `client`)
- `--client-name`, `-c`: Name for the generated client struct (default: `APIClient`)
- `--log-level`: Log level for output (debug, info, warn, error) (default: `info`)
- `--debug`: Enable debug logging (equivalent to `--log-level=debug`)
- `--json-log`: Output logs in JSON format for structured logging

### Logging

The tool provides comprehensive logging for debugging and monitoring:

- **Info Level** (default): Shows major operation progress and completion status
- **Debug Level**: Detailed tracing of all operations, including parameter extraction, file operations, and template processing
- **Warn/Error Levels**: Important warnings and error conditions with context
- **JSON Format**: Structured logging suitable for log aggregation systems

#### Logging Examples

```bash
# Default info-level logging
go-api-gen --input petstore.yaml --output ./generated

# Debug logging for troubleshooting
go-api-gen --input petstore.yaml --output ./generated --debug

# JSON structured logging for monitoring
go-api-gen --input petstore.yaml --output ./generated --json-log

# Custom log level
go-api-gen --input petstore.yaml --output ./generated --log-level warn
```

#### Log Components

Each component logs with a distinct identifier for easy filtering:
- `cli`: Command-line interface operations
- `generator`: Overall code generation orchestration
- `parser`: OpenAPI specification parsing and validation
- `client-template`: Client code generation and operation extraction
- `models-template`: Model and type generation

## Generated Code Structure

The tool generates two main files:

### client.go
Contains the main HTTP client with:
- Client struct with configuration
- Constructor functions
- HTTP request methods for each API endpoint
- Error handling and response parsing
- URL building with path and query parameters

### models.go
Contains data models with:
- Go structs for all schema definitions
- Proper JSON tags for serialization
- Type aliases for array types
- Time imports only when needed

## Example

Given this OpenAPI specification:

```yaml
openapi: 3.0.0
info:
  title: Pet Store API
  version: 1.0.0
servers:
  - url: https://petstore.example.com/v1
paths:
  /pets:
    get:
      operationId: listPets
      parameters:
        - name: limit
          in: query
          schema:
            type: integer
      responses:
        '200':
          description: A list of pets
components:
  schemas:
    Pet:
      type: object
      properties:
        id:
          type: integer
        name:
          type: string
```

The tool generates a client that can be used like:

```go
package main

import (
    "context"
    "log"
    
    "your-module/generated/client"
)

func main() {
    client := client.NewAPIClient()
    
    ctx := context.Background()
    limit := 10
    
    if err := client.ListPets(ctx, &limit); err != nil {
        log.Fatal(err)
    }
}
```

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build ./cmd/go-api-gen
```

## License

MIT License - see LICENSE file for details.