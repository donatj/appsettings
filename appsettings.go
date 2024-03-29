// Package appsettings provides simple key/value store functionality designed
// to be used for easily storing and recalling application runtime settings
package appsettings

import (
	"encoding/json"
	"errors"
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
	IncrInt64(key string) int64
	DecrInt64(key string) int64

	Delete(key string)
	DeleteTree(key string)
	GetTree(key string) DataTree

	HasTree(key string) bool
	HasLeaf(key string) bool

	GetTrees() map[string]DataTree
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
	pretty   bool

	storage StorageAdapter

	*tree
}

// Option sets an option of the passed AppSettings
type Option func(*AppSettings)

// OptionPrettyPrint configures AppSettings to pretty print the saved json
func OptionPrettyPrint(app *AppSettings) {
	app.pretty = true
}

// OptionStorageAdapter sets the storage adapter for the AppSettings
func OptionStorageAdapter(storage StorageAdapter) Option {
	return func(app *AppSettings) {
		app.storage = storage
	}
}

// NewAppSettings gets a new AppSettings struct
func NewAppSettings(dbFilename string, options ...Option) (*AppSettings, error) {
	a := &AppSettings{
		filename: dbFilename,
		pretty:   false,

		storage: NewFileSysPersister(),

		tree: &tree{
			Branches: make(map[string]*tree),
			Leaves:   make(map[string]string),
		},
	}

	for _, option := range options {
		option(a)
	}

	data, err := a.storage.Fetch(dbFilename)
	if errors.Is(err, ErrorEmptyFetch) {
		return a, nil
	} else if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// ErrUndefinedKey is returned when the key requested from get is undefined.
var ErrUndefinedKey = errors.New("undefined key")

// GetString gets the give keys leaf as a string.
//
// Returns an ErrUndefinedKey if the key is not defined
func (a *tree) GetString(key string) (string, error) {
	a.Lock()
	defer a.Unlock()

	if _, ok := a.Leaves[key]; !ok {
		return "", ErrUndefinedKey
	}

	return a.Leaves[key], nil
}

// SetString sets the give key leaf as a string.
func (a *tree) SetString(key string, val string) {
	a.Lock()
	defer a.Unlock()

	a.Leaves[key] = val
}

// GetInt gets the give keys leaf as an int.
//
// Returns an ErrUndefinedKey if the key is not defined.
//
// Will return an error if ParseInt of the value of the leaf fails.
func (a *tree) GetInt(key string) (int, error) {
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

// SetInt sets the give key leaf as a int.
func (a *tree) SetInt(key string, val int) {
	a.Lock()
	defer a.Unlock()

	a.Leaves[key] = strconv.Itoa(val)
}

// GetInt64 gets the give keys leaf as an int64.
//
// Returns an ErrUndefinedKey if the key is not defined.
//
// Will return an error if ParseInt(...,10, 64) of the value of the leaf fails.
func (a *tree) GetInt64(key string) (int64, error) {
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

// SetInt64 sets the give key leaf as a int64.
func (a *tree) SetInt64(key string, val int64) {
	a.Lock()
	defer a.Unlock()

	a.Leaves[key] = strconv.FormatInt(val, 10)
}

// IncrInt64 increments an int or int64 leaf and returns the new value as int64.
//
// If the key is undefined or non-integer, this will initialize it to 0 and increment it to 1
func (a *tree) IncrInt64(key string) int64 {
	v, _ := a.GetInt64(key)
	v += 1
	a.SetInt64(key, v)

	return v
}

// DecrInt64 decrements an int or int64 leaf and returns the new value as int64.
//
// If the key is undefined or non-integer, this will initialize it to 0 and decrements it to -1
func (a *tree) DecrInt64(key string) int64 {
	v, _ := a.GetInt64(key)
	v -= 1
	a.SetInt64(key, v)

	return v
}

// Delete removes the given leaf from the branch.
func (a *tree) Delete(key string) {
	a.Lock()
	defer a.Unlock()

	delete(a.Leaves, key)
}

// DeleteTree removes the given tree from the branch.
func (a *tree) DeleteTree(key string) {
	a.Lock()
	defer a.Unlock()

	delete(a.Branches, key)
}

func (a *tree) GetLeaves() map[string]string {
	return a.Leaves
}

func (a *tree) GetTrees() map[string]DataTree {
	out := make(map[string]DataTree)
	for k, b := range a.Branches {
		out[k] = b
	}

	return out
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

// HasTree checks if the given key is defined as a tree
func (a *tree) HasTree(key string) bool {
	a.Lock()
	defer a.Unlock()

	_, ok := a.Branches[key]
	return ok
}

// HasLeaf checks if the given key is defined as a leaf value
func (a *tree) HasLeaf(key string) bool {
	a.Lock()
	defer a.Unlock()

	_, ok := a.Leaves[key]
	return ok
}

// Persist causes the current state of the app settings to be persisted.
func (a *AppSettings) Persist() error {
	a.Lock()
	defer a.Unlock()

	var err error
	var d1 []byte

	if a.pretty {
		d1, err = json.MarshalIndent(a, "", "\t")
	} else {
		d1, err = json.Marshal(a)
	}
	if err != nil {
		return err
	}

	return a.storage.Persist(a.filename, d1)
}
