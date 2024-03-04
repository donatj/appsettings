package appsettings

import (
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
	Fetch(filename string) ([]byte, error)

	// Persist causes the current state of the app settings to be persisted.
	Persist(filename string, data []byte) error
}

// FileSystemStorageAdapter is a StorageAdapter that uses the file system to persist and
// fetch app settings. It is the default StorageAdapter for AppSettings.
type FileSystemStorageAdapter struct {
	sync.Mutex
}

var _ StorageAdapter = &FileSystemStorageAdapter{}

// NewFileSysPersister creates a new FileSystemStorageAdapter with the given filename.
func NewFileSysPersister() *FileSystemStorageAdapter {
	return &FileSystemStorageAdapter{}
}

func (f *FileSystemStorageAdapter) Persist(filename string, data []byte) error {
	f.Lock()
	defer f.Unlock()

	return f.writeFile(filename, data)
}

func (f *FileSystemStorageAdapter) writeFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}

// ErrorEmptyFetch is returned when the fetch result is empty.
var ErrorEmptyFetch = errors.New("empty fetch result")

// Fetch retrieves the persisted app settings from the file system.
func (f *FileSystemStorageAdapter) Fetch(filename string) ([]byte, error) {
	f.Lock()
	defer f.Unlock()

	var empty = []byte("{}")

	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		err := f.writeFile(filename, empty)
		if err != nil {
			return nil, err
		}

		return empty, ErrorEmptyFetch
	} else if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 || slices.Equal(data, empty) {
		return empty, ErrorEmptyFetch
	}

	return data, nil
}
