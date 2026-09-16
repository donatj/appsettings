# appsettings

[![CI](https://github.com/donatj/appsettings/actions/workflows/ci.yml/badge.svg)](https://github.com/donatj/appsettings/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/donatj/appsettings.svg)](https://pkg.go.dev/github.com/donatj/appsettings)

`appsettings` is a small Go package for persisting application runtime settings in a JSON-backed, hierarchical key/value store. It is useful when an application needs a little durable local state without defining a separate configuration type for every value.

- Organize values into named trees.
- Store and retrieve strings, `int`s, and `int64`s.
- Use the included JSON file storage or provide a storage adapter of your own.
- Manage the same data from a small command-line tool.

## Install

`appsettings` requires Go 1.21 or newer.

```sh
go get github.com/donatj/appsettings
```

## Basic usage

Create an `AppSettings` value with a filename, work with a tree, and call `Persist` when the changes should be written. With the default storage adapter, a missing file is created as an empty JSON object.

```go
package main

import (
	"errors"
	"log"

	"github.com/donatj/appsettings"
)

func main() {
	settings, err := appsettings.NewAppSettings(
		"appsettings.json",
		appsettings.OptionPrettyPrint,
	)
	if err != nil {
		log.Fatal(err)
	}

	user := settings.GetTree("user")

	if theme, err := user.GetString("theme"); err == nil {
		log.Printf("using %s theme", theme)
	} else if errors.Is(err, appsettings.ErrUndefinedKey) {
		user.SetString("theme", "system")
	} else {
		log.Fatal(err)
	}

	user.SetInt("window-width", 1280)
	log.Printf("launch %d", user.IncrInt64("launch-count"))

	if err := settings.Persist(); err != nil {
		log.Fatal(err)
	}
}
```

`GetTree` returns a child tree and creates it when it does not already exist. Leaves are stored as strings; the integer setters convert their values to decimal strings and the integer getters parse them. `IncrInt64` and `DecrInt64` return the updated value; an undefined or non-integer value is initialized before changing it.

Call `HasLeaf` or `HasTree` to check for a value or child tree without creating anything. Use `Delete` to remove a leaf and `DeleteTree` to remove a child tree.

## Persistence

The default `FileSystemStorageAdapter` reads and writes the filename supplied to `NewAppSettings`. Changes stay in memory until `Persist` succeeds. Pass `OptionPrettyPrint` to write indented JSON instead of compact JSON.

For another backing store, implement `StorageAdapter` and pass it with `OptionStorageAdapter`:

```go
settings, err := appsettings.NewAppSettings(
	"settings-for-current-user",
	appsettings.OptionStorageAdapter(adapter),
)
```

An adapter's `Fetch` method should return `appsettings.ErrorEmptyFetch` when there is no saved state. This lets `NewAppSettings` return an initialized, empty store. Its `Persist` method receives the same name and the serialized JSON to save.

## Command-line tool

Install the CLI with:

```sh
go install github.com/donatj/appsettings/cmd/appsettings@latest
```

The `-file` flag selects the JSON file and must appear before the subcommand. Paths use dots to traverse trees.

```sh
# Set one or more path/value pairs.
appsettings -file ./settings.json set user.name Jesse user.theme dark

# Print one or more values, one per line.
appsettings -file ./settings.json get user.name user.theme

# Remove one or more leaf values.
appsettings -file ./settings.json delete user.theme
```

The CLI stores values as strings. `set` requires an even number of path/value arguments; `get` and `delete` operate on each path supplied.

## Documentation and license

See the [package documentation](https://pkg.go.dev/github.com/donatj/appsettings) for the complete API. This project is released under the [MIT License](LICENSE.md).
