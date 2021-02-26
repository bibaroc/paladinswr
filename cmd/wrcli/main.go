package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bibaroc/paladinswr/stats"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

func main() {
	inflixDBConfig := getWriteAPIConfig()
	client := influxdb2.NewClient(inflixDBConfig.url, inflixDBConfig.token)

	defer client.Close()

	writeAPI := client.WriteAPI(inflixDBConfig.org, inflixDBConfig.bucket)
	go func(writer api.WriteAPI) {
		for err := range writer.Errors() {
			log.Println(err.Error())
		}
	}(writeAPI)

	stats, err := stats.ListStats()
	if err != nil {
		log.Println(err)
		return
	}

	for _, st := range stats {
		p := influxdb2.NewPointWithMeasurement("wr").
			AddTag("unit", "percentage").
			AddTag("class", st.Class).
			AddTag("champion", st.Name).
			AddField("max", st.Winrate.BestLoadout).
			AddField("med", st.Winrate.WeightedAverage).
			SetTime(time.Now())
		writeAPI.WritePoint(p)
	}

	writeAPI.Flush()
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
