package appsettings

import (
	"encoding/json"
	"errors"
	"os"
	"slices"
	"sync"
)

// StorageAdapter is an interface for persisting and fetching serialized app settings
type StorageAdapter interface {
	// Fetch retrieves the persisted app settings.
	// If the result is empty, it will return ErrorEmptyFetch
	// Any initialization should be done in the Fetch method.
	Fetch() ([]byte, error)

	// Persist causes the current state of the app settings to be persisted.
	Persist([]byte) error
}

// FileSystemStorageAdapter is a StorageAdapter that uses the file system to persist and
// fetch app settings. It is the default StorageAdapter for AppSettings.
type FileSystemStorageAdapter struct {
	filename string
	sync.Mutex
}

var _ StorageAdapter = &FileSystemStorageAdapter{}

// NewFileSysPersister creates a new FileSystemStorageAdapter with the given filename.
func NewFileSysPersister(filename string) *FileSystemStorageAdapter {
	return &FileSystemStorageAdapter{filename: filename}
}

func (f *FileSystemStorageAdapter) Persist(data []byte) error {
	f.Lock()
	defer f.Unlock()

	return f.writeFile(data)
}

func (f *FileSystemStorageAdapter) writeFile(data []byte) error {
	return os.WriteFile(f.filename, data, 0644)
}

// ErrorEmptyFetch is returned when the fetch result is empty.
var ErrorEmptyFetch = errors.New("empty fetch result")

// Fetch retrieves the persisted app settings from the file system.
func (f *FileSystemStorageAdapter) Fetch() ([]byte, error) {
	f.Lock()
	defer f.Unlock()

	var empty = []byte("{}")

	_, err := os.Stat(f.filename)
	if os.IsNotExist(err) {
		err := f.writeFile(empty)
		if err != nil {
			return nil, err
		}

		return empty, ErrorEmptyFetch
	} else if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(f.filename)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 || slices.Equal(data, empty) {
		return empty, ErrorEmptyFetch
	}

	return data, nil
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

	return a.storage.Persist(d1)
}
