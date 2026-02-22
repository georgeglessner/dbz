# AGENTS.md - Coding Guidelines for DBZ

## Build Commands

```bash
make build          # Build the binary
make ci-build       # Build with version info
make release        # Cross-compile for all platforms
make clean          # Clean build artifacts
make install        # Install to /usr/local/bin
```

## Test Commands

```bash
make test                                   # Run all tests
go test -v ./pkg/containers                 # Run tests in a package
go test -v ./pkg/containers -run TestName   # Run a specific test
make ci-test                                # Run tests with coverage
make coverage                               # Generate coverage report
```

## Lint/Format Commands

```bash
make fmt    # Format code
make lint   # Run linter
make vet    # Run go vet
make deps   # Download and tidy dependencies
```

## Code Style Guidelines

### Imports
Group imports: stdlib first, then third-party, then local. Use `goimports` format.

```go
import (
    "context"
    "fmt"
    
    "github.com/docker/docker/api/types"
    "github.com/spf13/cobra"
    
    "github.com/dbz/dbz/pkg/containers"
)
```

### Formatting
- Use `gofmt` for all code
- Line length: prefer under 100 characters
- Use tabs for indentation (Go standard)
- No trailing whitespace

### Naming Conventions
- **Types/Interfaces**: PascalCase (e.g., `Database`, `ContainerInfo`)
- **Functions**: PascalCase exported, camelCase unexported
- **Variables**: camelCase (e.g., `containerName`, `dbType`)
- **Constants**: PascalCase for exported
- **Interface names**: Noun-based (e.g., `Database`)
- **Test functions**: `Test<Name>` with table-driven tests
- **Packages**: lowercase, single word (e.g., `containers`, `cmd`)

### Types and Structs
- Define interfaces in `types.go` files
- Prefer composition over inheritance
- Document exported types with comments starting with the type name

```go
// ContainerConfig holds configuration for creating a database container
type ContainerConfig struct {
    Type     string
    Port     int
    Password string
}
```

### Error Handling
- Always check errors and return them wrapped with context
- Use `fmt.Errorf` with `%w` verb for error wrapping
- Return errors, don't log and continue (except cleanup)
- Error messages: lowercase, no punctuation

```go
if err != nil {
    return nil, fmt.Errorf("failed to create container: %w", err)
}
```

### Functions
- Keep functions focused and under 50 lines
- Return early to reduce nesting
- Context should be first parameter
- Error should be last return value

### Comments
- All exported functions, types, and constants must have doc comments
- Comments start with the name of the thing being documented
- Use complete sentences with proper punctuation

### Testing
- Use table-driven tests with `tests := []struct{...}`
- Test names: descriptive of what's being tested
- Use `t.Run(tt.name, func(t *testing.T) {...})` for subtests
- Check both error and non-error cases

```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid case", "input", "output", false},
        {"invalid case", "", "", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Function() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("Function() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### CLI Commands (Cobra)
- Use `cmd/` package for command definitions
- Define flags in `init()` with `Flags()`
- Use `RunE` to return errors (not `Run`)
- Provide both short and long descriptions
- Include examples in command help

### Architecture Patterns
- Use factory pattern for database type creation (`DatabaseFactory`)
- Implement interface-based design (`Database` interface)
- Keep Docker/container logic in `pkg/containers`
- Keep CLI commands in `cmd/`
- Use `pkg/` for library code

### Dependencies
- Minimize external dependencies
- Prefer standard library when possible
- Run `go mod tidy` after adding/removing imports

## Running the Application

```bash
make run              # Build and run
make dev              # Development mode
make create-postgres  # Create test databases
make create-mysql
make list             # List databases
make cleanup          # Clean up test databases
```
