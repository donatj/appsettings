// Package appsettings provides simple key/value store functionality
package appsettings

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
)

// DataTree is the host of key and branches of values
type DataTree interface {
	GetString(key string) (string, error)
	SetString(key string, val string)
	GetInt(key string) (int, error)
	SetInt(key string, val int)
	GetInt64(key string) (int64, error)
	SetInt64(key string, val int64)
	Delete(key string)
	GetTree(key string) DataTree

	GetLeaves() map[string]string
}

type tree struct {
	Branches map[string]*tree
	Leaves   map[string]string

	sync.Mutex
}

// AppSettings is the root most DataTree
type AppSettings struct {
	filename string

	*tree
}

// NewAppSettings gets a new AppSettings struct
func NewAppSettings(dbFilename string) (*AppSettings, error) {
	a := &AppSettings{
		filename: dbFilename,
		tree: &tree{
			Branches: make(map[string]*tree),
			Leaves:   make(map[string]string),
		},
	}

	if _, err := os.Stat(dbFilename); os.IsNotExist(err) {
		d1, _ := json.Marshal(a)
		err := ioutil.WriteFile(dbFilename, d1, 0644)
		if err != nil {
			return nil, err
		}
	} else {
		d1, err := ioutil.ReadFile(dbFilename)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(d1, a)
		if err != nil {
			return nil, err
		}
	}

	return a, nil
}

// ErrUndefinedKey is returned when the key requested from get is undefined.
var ErrUndefinedKey = errors.New("undefined key")

func (a *tree) GetString(key string) (string, error) {
	a.Lock()
	defer a.Unlock()

	if _, ok := a.Leaves[key]; !ok {
		return "", ErrUndefinedKey
	}

	return a.Leaves[key], nil
}

func (a *tree) SetString(key string, val string) {
	a.Lock()
	defer a.Unlock()

	a.Leaves[key] = val
}

func (a *tree) GetInt(key string) (int, error) {
	a.Lock()
	defer a.Unlock()

	str, err := a.GetString(key)
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (a *tree) SetInt(key string, val int) {
	a.Lock()
	defer a.Unlock()

	a.Leaves[key] = strconv.Itoa(val)
}

func (a *tree) GetInt64(key string) (int64, error) {
	a.Lock()
	defer a.Unlock()

	str, err := a.GetString(key)
	if err != nil {
		return 0, err
	}

	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (a *tree) SetInt64(key string, val int64) {
	a.Lock()
	defer a.Unlock()

	a.Leaves[key] = strconv.FormatInt(val, 10)
}

func (a *tree) Delete(key string) {
	a.Lock()
	defer a.Unlock()

	delete(a.Leaves, key)
}

func (a *tree) GetLeaves() map[string]string {
	return a.Leaves
}

// GetTree fetches a tree for app setting storage
func (a *tree) GetTree(key string) DataTree {
	a.Lock()
	defer a.Unlock()

	if _, ok := a.Branches[key]; !ok {
		a.Branches[key] = &tree{
			Branches: make(map[string]*tree),
			Leaves:   make(map[string]string),
		}
	}

	return a.Branches[key]
}

// Persist causes the current state of the app settings to be persisted.
func (a *AppSettings) Persist() error {
	a.Lock()
	defer a.Unlock()

	d1, err := json.Marshal(a)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(a.filename, d1, 0644)
	if err != nil {
		return err
	}

	return nil
}
