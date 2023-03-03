package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/open-policy-agent/opa/rego"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

type Data struct {
	ApiKeys map[string]interface{} `json:"ApiKeys"`
}

func main() {
	println("Loading API Keys into memory...")
	base := rego.New(rego.Query("data.apikeys.allow"),
		rego.LoadBundle("./bundle/bundle.tar.gz"))

	preparedQuery, err := base.PrepareForEval(context.Background())
	if err != nil {
		panic(err)
	}

	file, _ := ioutil.ReadFile("./bundle/apikeys/data.json")
	data := Data{}
	_ = json.Unmarshal([]byte(file), &data)

	allApiKeys := make([]string, 0, 1000000)

	for k, _ := range data.ApiKeys {
		allApiKeys = append(allApiKeys, k)
	}

	println("API Keys loaded into memory.")

	differentInputs := make([]interface{}, 0, 100)

	// Add 50 existing api keys picked randomly from data.json
	for i := 0; i < 50; i++ {
		randKey := allApiKeys[rand.Intn(1000000)]
		inputRaw := fmt.Sprintf("{\"apikey\": \"%s\"}`", randKey)
		var input interface{}
		err = json.NewDecoder(bytes.NewBufferString(inputRaw)).Decode(&input)
		if err != nil {
			panic(err)
		}
		differentInputs = append(differentInputs, input)
	}

	// Add 50 non-existing api keys
	for i := 0; i < 50; i++ {
		guuid := uuid.NewV4()
		inputRaw := fmt.Sprintf("{\"apikey\": \"%s\"}`", guuid.String())
		var input interface{}
		err = json.NewDecoder(bytes.NewBufferString(inputRaw)).Decode(&input)
		if err != nil {
			panic(err)
		}
		differentInputs = append(differentInputs, input)
	}

	// Test the 100 api keys
	allowedCount := 0
	rejectedCount := 0
	start := time.Now()
	for _, input := range differentInputs {
		result, err := preparedQuery.Eval(context.Background(), rego.EvalInput(input))
		if err != nil {
			panic(err)
		}
		if result.Allowed() {
			allowedCount++
		} else {
			rejectedCount++
		}
	}

	elapsed := time.Since(start)
	println()
	println("Allowed count: ", allowedCount)
	println("Rejected count: ", rejectedCount)
	log.Printf("Avg evaluation took: %s", (elapsed / 100))
}
