# udotEnv

`udotEnv` is a Go package for managing environment variables in your application. It provides a flexible configuration system for loading `.env` files and supports custom flags for environment variables and overload options.

## Features

- Load environment variables from a `.env` file.
- Support for custom flags to specify environment files and overload options.
- Default configuration with predefined flags and file paths.
- Panic handling for invalid configurations or duplicate flags.

## Installation

To use `udotEnv` in your project, add it to your `go.mod` file:

```bash
go get github.com/kravlad/udotEnv
```

## Usage

### Basic Usage

To use `udotEnvType` with the default configuration:

```go
package main

import (
    "github.com/kravlad/go-udotenv"
)

func main() {
    udotEnv := udotEnv.New(true).Load()
}
```

### Custom Configuration

You can provide a custom configuration to override the default settings:

```go
package main

import (
    "github.com/kravlad/go-udotenv"
)

func main() {
    customConfig := &udotEnv.Config{
        EnvFlags:       []string{"my-env", "E"},
        OverloadFlags:  []string{"custom-overload"},
        DefaultEnvPath: ".env.custom",
    }
    udotEnv := udotEnv.New(true, customConfig)
    udotEnv.Load()
}
```

### Handling Flags

`udotEnv` allows you to specify flags for environment files and overload options. For example:

```bash
./your-app --envs .env.test --env-overload --envs .env
```

### Default Configuration

The default configuration includes:

- `EnvFlags`: `["envs", "e"]`
- `OverloadFlags`: `["env-overload", "eo", "o"]`
- `DefaultEnvPath`: `.env`

### Example

```go
package main

import (
    "github.com/kravlad/go-udotenv"
)

func main() {
    udotEnv := udotEnv.New(true)
    udotEnv.Load()
}
```

## API Reference

### `func GetDefaultConfig() *Config`

Returns a pointer to a `Config` struct initialized with default values.

### `func New(parseFlags bool, config ...*Config) *udotEnvType`

Creates and initializes a new instance of `udotEnvType`.

- `parseFlags`: Whether to parse command-line flags immediately.
- `config`: Optional custom configuration.

### `func (ue *udotEnvType) Load()`

Loads environment variables from the specified file.

## Testing

Run the tests using the `go test` command:

```bash
go test
```

## License

This project is licensed under the MIT License.
