// Package appsettings provides simple key/value store functionality
package appsettings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
)

type DataTree map[string]string

type dataStruct struct {
	Tree map[string]DataTree
}

type AppSettings struct {
	filename string
	data     dataStruct

	sync.Mutex
}

func NewAppSettings(dbFilename string) (*AppSettings, error) {
	var data dataStruct
	if _, err := os.Stat(dbFilename); os.IsNotExist(err) {
		data = dataStruct{
			Tree: make(map[string]DataTree),
		}

		d1, _ := json.Marshal(data)
		err := ioutil.WriteFile(dbFilename, d1, 0644)
		if err != nil {
			return nil, err
		}
	} else {
		d1, err := ioutil.ReadFile(dbFilename)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(d1, &data)
		if err != nil {
			return nil, err
		}
	}

	return &AppSettings{filename: dbFilename, data: data}, nil
}

func (a *DataTree) GetString(key string) (string, error) {
	y := *a
	if _, ok := y[key]; !ok {
		return "", fmt.Errorf("undefined key %s", key)
	}

	return y[key], nil
}

func (a *DataTree) SetString(key string, val string) {
	y := *a
	y[key] = val
}

func (a *DataTree) GetInt(key string) (int, error) {
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

func (a *DataTree) SetInt(key string, val int) {
	y := *a
	y[key] = strconv.Itoa(val)
}

func (a *DataTree) GetInt64(key string) (int64, error) {
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

func (a *DataTree) SetInt64(key string, val int64) {
	y := *a
	y[key] = strconv.FormatInt(val, 10)
}

func (a *DataTree) Delete(key string) {
	delete(*a, key)
}

func (a *AppSettings) GetTree(key string) DataTree {
	a.Lock()
	defer a.Unlock()

	if _, ok := a.data.Tree[key]; !ok {
		a.data.Tree[key] = make(DataTree)
	}

	return a.data.Tree[key]
}

// Persist causes the current state of the app settings to be persisted.
func (a *AppSettings) Persist() error {
	a.Lock()
	defer a.Unlock()

	d1, _ := json.Marshal(a.data)
	err := ioutil.WriteFile(a.filename, d1, 0644)
	if err != nil {
		return err
	}

	return nil
}
