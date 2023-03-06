package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/open-policy-agent/opa/sdk"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"
)

type Data struct {
	ApiKeys map[string]interface{} `json:"ApiKeys"`
}

var (
	api_keys_count      = 1000000
	total_keys_to_check = 100
)

func main() {
	ctx := context.Background()
	println("Loading API Keys into memory...")

	opa, err := sdk.New(ctx, sdk.Options{
		Config: strings.NewReader(
			`
services:
  - name: hostedbundle2
    url: http://0.0.0.0:8081/

bundles:
  authz:
    service: hostedbundle2
    resource: bundle.tar.gz
    polling:
      min_delay_seconds: 60
      max_delay_seconds: 120
`)})
	if err != nil {
		panic(err)
	}

	file, err := ioutil.ReadFile("bundle/apikeys/data.json")
	if err != nil {
		panic(err)
	}

	data := Data{}
	_ = json.Unmarshal([]byte(file), &data)

	allApiKeys := make([]string, 0, api_keys_count)

	for k, _ := range data.ApiKeys {
		allApiKeys = append(allApiKeys, k)
	}

	println("API Keys loaded into memory.")

	differentInputs := make([]interface{}, 0, total_keys_to_check)

	// Add 50 existing api keys picked randomly from data.json
	for i := 0; i < total_keys_to_check/2; i++ {
		randKey := allApiKeys[rand.Intn(api_keys_count)]
		inputRaw := fmt.Sprintf("{\"apikey\": \"%s\"}`", randKey)
		var input interface{}
		err = json.NewDecoder(bytes.NewBufferString(inputRaw)).Decode(&input)
		if err != nil {
			panic(err)
		}
		differentInputs = append(differentInputs, input)
	}

	// Add 50 non-existing api keys
	for i := 0; i < total_keys_to_check/2; i++ {
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
		decision, err := opa.Decision(ctx, sdk.DecisionOptions{Path: "/apikeys/allow", Input: input})
		if err != nil {
			println(err.Error())
		}
		if decision.Result != nil && decision.Result.(bool) {
			allowedCount++
		} else {
			rejectedCount++
		}
	}

	elapsed := time.Since(start)
	println()
	println("Allowed count: ", allowedCount)
	println("Rejected count: ", rejectedCount)
	log.Printf("Avg evaluation took: %s", (elapsed.Microseconds() / 100))
}
