package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/open-policy-agent/opa/rego"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"math/rand"
	"runtime"
	"time"
)

type Data struct {
	ApiKeys map[string]interface{} `json:"ApiKeys"`
}

const totalKeysToCheck = 100
const halfKeysToCheck = totalKeysToCheck / 2

func main() {
	fmt.Println("Loading API Keys into memory...")

	preparedQuery := prepareRegoQuery()
	existingKeys := getExistingKeys()

	differentInputs := appendInvalidKeysAndShuffle(existingKeys)

	start := time.Now()
	allowedCount, rejectedCount := evaluateApiKeys(preparedQuery, differentInputs)
	elapsedTime := time.Since(start).Microseconds()

	fmt.Printf("\nAllowed count: %d\n", allowedCount)
	fmt.Printf("Rejected count: %d\n", rejectedCount)
	fmt.Printf("Avg evaluation took: %dms\n", elapsedTime/totalKeysToCheck)

	runtime.GC()
	printMemUsage()
}

func prepareRegoQuery() *rego.PreparedEvalQuery {
	base := rego.New(rego.Query("data.apikeys.allow"), rego.LoadBundle("./bundle/bundle.tar.gz"))
	preparedQuery, err := base.PrepareForEval(context.Background())
	if err != nil {
		log.Fatalf("Failed preparing rego query: %v", err)
	}
	return &preparedQuery
}

func getExistingKeys() []string {
	file, err := ioutil.ReadFile("./bundle/apikeys/data.json")
	if err != nil {
		log.Fatalf("Failed reading file: %v", err)
	}

	data := Data{}
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
		randKey := allApiKeys[rand.Intn(halfKeysToCheck)]
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

func evaluateApiKeys(preparedQuery *rego.PreparedEvalQuery, differentInputs []interface{}) (int, int) {
	allowedCount := 0
	rejectedCount := 0
	for _, input := range differentInputs {
		result, err := preparedQuery.Eval(context.Background(), rego.EvalInput(input))
		if err != nil {
			log.Println("Error in decision:", err)
			continue
		}
		if result.Allowed() {
			allowedCount++
		} else {
			rejectedCount++
		}
	}
	return allowedCount, rejectedCount
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
