package main

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"io/ioutil"
)

type Data struct {
	ApiKeys map[string]interface{} `json:"ApiKeys"`
}

func main() {
	data := Data{
		ApiKeys: make(map[string]interface{}),
	}
	for i := 0; i < 1000000; i++ {
		// generate random uuid
		guuid := uuid.NewV4()
		data.ApiKeys[guuid.String()] = struct{}{}
	}

	file, _ := json.MarshalIndent(data, "", " ")

	_ = ioutil.WriteFile("./bundle/apikeys/data.json", file, 0644)
}
