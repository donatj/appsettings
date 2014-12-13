# AppSettings

[![GoDoc](https://godoc.org/github.com/donatj/appsettings?status.svg)](https://godoc.org/github.com/donatj/appsettings)

A dead simple key value store for persisting simple runtime options in Go applications.

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

## Installation

```
go get github.com/donatj/appsettings
```

## Documentation

Documentation can be found a godoc:

https://godoc.org/github.com/donatj/appsettings
