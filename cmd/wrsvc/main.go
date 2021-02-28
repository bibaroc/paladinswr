package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	inflixDBConfig := getWriteAPIConfig()
	client := influxdb2.NewClient(inflixDBConfig.url, inflixDBConfig.token)

	defer client.Close()
	query := fmt.Sprintf("from(bucket:\"%v\")|> range(start: -30d) |> filter(fn: (r) => r._measurement == \"wr\")", inflixDBConfig.bucket)
	queryAPI := client.QueryAPI(inflixDBConfig.org)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		log.Println(err)
		return
	}

	response := make(map[string]map[string]map[time.Time]StatPoint)

	for result.Next() {
		championsClass := result.Record().ValueByKey("class").(string)
		championsName := result.Record().ValueByKey("champion").(string)
		statTime := result.Record().Time()

		if _, ok := response[championsClass]; !ok {
			fmt.Println("response[championsClass] is nill")
			response[championsClass] = map[string]map[time.Time]StatPoint{}
		}

		if _, ok := response[championsClass][championsName]; !ok {
			fmt.Println("response[championsClass][championsName] is nill")
			response[championsClass][championsName] = map[time.Time]StatPoint{}
		}

		switch result.Record().Field() {
		case "max":
			response[championsClass][championsName][statTime] = StatPoint{
				Average: response[championsClass][championsName][statTime].Average,
				Max:     result.Record().Value().(float64),
			}
		case "avg":
			response[championsClass][championsName][statTime] = StatPoint{
				Average: result.Record().Value().(float64),
				Max:     response[championsClass][championsName][statTime].Max,
			}
		}
	}
	// check for an error
	if result.Err() != nil {
		fmt.Printf("query parsing error: %s\n", result.Err().Error())
	}

	formatted, _ := json.MarshalIndent(response, "", "\t")
	fmt.Println(string(formatted))
}

type writeAPIConfg struct {
	token  string
	bucket string
	org    string
	url    string
}

func getWriteAPIConfig() writeAPIConfg {
	mustString := func(s string) string {
		if v, ok := os.LookupEnv(s); ok {
			return v
		}

		panic(fmt.Sprintf("environment value for %q not found", s))
	}

	return writeAPIConfg{
		token:  mustString("WRCLI_TOKEN"),
		bucket: mustString("WRCLI_BUCKET"),
		org:    mustString("WRCLI_ORG"),
		url:    mustString("WRCLI_URL"),
	}
}

type StatPoint struct {
	Average float64 `json:"avg,omitempty"`
	Max     float64 `json:"max,omitempty"`
}
