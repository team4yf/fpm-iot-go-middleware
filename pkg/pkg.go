package main

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/oliveagle/jsonpath"
)

func GetBodyString(readCloser io.ReadCloser) (string, error) {
	body, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return "", err
	}
	defer readCloser.Close()

	return string(body), nil
}

func GetBodyMap(readCloser io.ReadCloser) (map[string]interface{}, error) {
	var body map[string]interface{}
	err := json.NewDecoder(readCloser).Decode(&body)
	if err != nil {
		return nil, err
	}
	defer readCloser.Close()

	return body, nil
}

func GetJsonPathData(data, jp string) (interface{}, error) {
	var jsonData interface{}
	json.Unmarshal([]byte(data), &jsonData)
	res, err := jsonpath.JsonPathLookup(jsonData, jp)
	if err != nil {
		return nil, err
	}
	return res, nil
}
