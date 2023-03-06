package main

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"math/rand"
)

type Data struct {
	ApiKeys map[string]interface{} `json:"ApiKeys"`
}

func main() {
	data := Data{
		ApiKeys: make(map[string]interface{}),
	}
	api_keys := 1000000
	for i := 0; i < api_keys; i++ {
		// generate random uuid
		guuid := uuid.NewV4()

		secretData := map[string]string{
			"Key":             randSeq(128),
			"UsagePlanName":   randSeq(68),
			"Username":        randSeq(32),
			"Email":           randSeq(32),
			"ApiProductName":  randSeq(68),
			"EnvironmentName": randSeq(68),
		}

		data.ApiKeys[guuid.String()] = secretData
	}

	file, _ := json.MarshalIndent(data, "", " ")

	_ = ioutil.WriteFile("./bundle/apikeys/data.json", file, 0644)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
