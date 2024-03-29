package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/open-policy-agent/opa/sdk"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

type Data struct {
	ApiKeys map[string]interface{} `json:"ApiKeys"`
}

const (
	apiKeysCount     = 1000000
	totalKeysToCheck = 100
	halfKeysToCheck  = totalKeysToCheck / 2
)

func main() {
	ctx := context.Background()
	fmt.Println("Loading API Keys into memory...")

	opa, err := sdk.New(ctx, sdk.Options{
		Config: bytes.NewReader([]byte(`
services:
  - name: hostedbundle2
    url: http://0.0.0.0:8981/

bundles:
  authz:
    service: hostedbundle2
    resource: bundle.tar.gz
    polling:
      min_delay_seconds: 60
      max_delay_seconds: 120
`))})
	if err != nil {
		log.Fatalf("Failed initializing SDK: %v", err)
	}

	existingKeys := getExistingKeys()
	differentInputs := appendInvalidKeysAndShuffle(existingKeys)

	start := time.Now()
	fmt.Println("API Keys loaded into memory.")
	// Test the 100 api keys, 50 should be accepted, and 50 should be rejected
	allowedCount, rejectedCount := evaluateApiKeys(opa, ctx, differentInputs)
	elapsedTime := time.Since(start).Microseconds()
	fmt.Printf("\nAllowed count: %d\n", allowedCount)
	fmt.Printf("Rejected count: %d\n", rejectedCount)
	fmt.Printf("Avg evaluation took: %d microseconds\n", elapsedTime/totalKeysToCheck)
}

func getExistingKeys() []string {
	// Read data in order to have valid inputs
	// So that we include in performance evaluation successfull cases and failures
	file, err := ioutil.ReadFile("bundle/apikeys/data.json")
	if err != nil {
		log.Fatalf("Failed reading file: %v", err)
	}

	var data Data
	if err := json.Unmarshal(file, &data); err != nil {
		log.Fatalf("Failed unmarshalling data: %v", err)
	}

	allApiKeys := make([]string, len(data.ApiKeys))
	i := 0
	for k := range data.ApiKeys {
		allApiKeys[i] = k
		i++
	}
	return allApiKeys
}

func appendInvalidKeysAndShuffle(allApiKeys []string) []interface{} {
	var differentInputs []interface{}
	for i := 0; i < halfKeysToCheck; i++ {
		// Add existing api keys picked randomly from data.json
		randKey := allApiKeys[rand.Intn(apiKeysCount/2)]
		differentInputs = append(differentInputs, createInput(randKey))

		// Add non-existing api keys
		guuid := uuid.NewV4()
		differentInputs = append(differentInputs, createInput(guuid.String()))
	}
	return differentInputs
}

func createInput(key string) interface{} {
	inputRaw := fmt.Sprintf("{\"apikey\": \"%s\"}", key)
	var input interface{}
	if err := json.NewDecoder(bytes.NewBufferString(inputRaw)).Decode(&input); err != nil {
		log.Fatalf("Failed decoding input: %v", err)
	}
	return input
}

func evaluateApiKeys(opa *sdk.OPA, ctx context.Context, differentInputs []interface{}) (int, int) {
	allowedCount := 0
	rejectedCount := 0
	for _, input := range differentInputs {
		decision, err := opa.Decision(ctx, sdk.DecisionOptions{Path: "/apikeys/allow", Input: input})
		if err != nil {
			log.Println("Error in decision:", err)
			continue
		}
		if result, ok := decision.Result.(bool); ok && result {
			allowedCount++
		} else {
			rejectedCount++
		}
	}
	return allowedCount, rejectedCount
}
