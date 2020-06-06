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
	if err := json.NewDecoder(readCloser).Decode(&body); err != nil {
		return nil, err
	}
	defer readCloser.Close()

	return body, nil
}

func GetJsonPathData(data, jp string) (interface{}, error) {
	var jsonData interface{}
	json.Unmarshal([]byte(data), &jsonData)
	if res, err := jsonpath.JsonPathLookup(jsonData, jp); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}
