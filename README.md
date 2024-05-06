# AppSettings

[![CI](https://github.com/donatj/appsettings/actions/workflows/ci.yml/badge.svg)](https://github.com/donatj/appsettings/actions/workflows/ci.yml)
[![GoDoc](https://godoc.org/github.com/donatj/appsettings?status.svg)](https://godoc.org/github.com/donatj/appsettings)

A hierarchical key value store for persisting simple runtime options in Go applications.

## Example

```go
s := appsettings.NewAppSettings("settings.json")

t := s.GetTree("user-settings")

//set
t.SetString("pizza", "pie")
t.SetInt("how-many-pugs", 349)

//read
if v, err := t.GetInt("how-many-pugs"); err == nil {
	log.Println(v)
}

if v, err = t.GetString("pizza"); err == nil {
	log.Println(v)
}

s.Persist()
```

## CLI Tool Installation

### From Source

```bash
$ go install github.com/donatj/appsettings/cmd/appsettings@latest
```

## Migration from v0.0.1

The JSON format for the early Alpha changed. To migrate your existing data compatibly to the more modern format, you can use [jq](https://stedolan.github.io/jq/) and execute the following command, first replacing `{your-file}` with the path to your actual database file.

```bash
jq '.Tree |= with_entries(.value = {Leaves: .value} ) | . + {Branches: .Tree} | del(.Tree)' < {your-file} > tmp && mv tmp {your-fie}
```

## Documentation

Documentation can be found a godoc:

https://godoc.org/github.com/donatj/appsettings
