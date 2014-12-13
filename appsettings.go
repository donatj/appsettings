package appsettings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type dataTree map[string]string

type dataStruct struct {
	Tree map[string]dataTree
}

type appsettings struct {
	filename string
	data     dataStruct
}

func NewAppSettings(dbFilename string) (*appsettings, error) {
	var data dataStruct
	if _, err := os.Stat(dbFilename); os.IsNotExist(err) {
		log.Printf("no such file or directory: %s", dbFilename)
		data = dataStruct{}

		d1, _ := json.Marshal(data)
		err := ioutil.WriteFile(dbFilename, d1, 0644)
		if err != nil {
			return nil, err
		}
	} else {
		d1, err := ioutil.ReadFile(dbFilename)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		json.Unmarshal(d1, &data)

		if err != nil {
			return nil, err
		}
	}

	return &appsettings{filename: dbFilename, data: data}, nil
}

func (a *dataTree) GetString(key string) (string, error) {
	y := *a
	if _, ok := y[key]; !ok {
		return "", fmt.Errorf("Undefined key %s", key)
	}

	return y[key], nil
}

func (a *dataTree) SetString(key string, val string) {
	y := *a
	y[key] = val
}

func (a *dataTree) GetInt(key string) (int, error) {
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

func (a *dataTree) SetInt(key string, val int) {
	y := *a
	y[key] = strconv.Itoa(val)
}

func (a *appsettings) GetTree(key string) dataTree {
	if _, ok := a.data.Tree[key]; !ok {
		a.data.Tree[key] = make(dataTree)
	}

	return a.data.Tree[key]
}

func (a *appsettings) Persist() error {
	d1, _ := json.Marshal(a.data)
	err := ioutil.WriteFile(a.filename, d1, 0644)
	if err != nil {
		return err
	}

	return nil
}
