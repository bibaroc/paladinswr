package wrsvc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-kit/kit/log"

	"net/http"
	"sync"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func GetAllChampionStatsProvider(
	client influxdb2.Client,
	bucket string,
	orgName string,
) func() (map[string]map[string]map[time.Time]StatPoint, error) {
	return func() (map[string]map[string]map[time.Time]StatPoint, error) {
		query := fmt.Sprintf("from(bucket:\"%v\")|> range(start: -30d) |> filter(fn: (r) => r._measurement == \"wr\")", bucket)
		queryAPI := client.QueryAPI(orgName)

		result, err := queryAPI.Query(context.Background(), query)
		if err != nil {
			return nil, fmt.Errorf("could not query data: %w", err)
		}

		response := make(map[string]map[string]map[time.Time]StatPoint)

		for result.Next() {
			championsClass := result.Record().ValueByKey("class").(string)
			championsName := result.Record().ValueByKey("champion").(string)
			statTime := result.Record().Time()

			if _, ok := response[championsClass]; !ok {
				response[championsClass] = map[string]map[time.Time]StatPoint{}
			}

			if _, ok := response[championsClass][championsName]; !ok {
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

		if err = result.Err(); err != nil {
			return nil, fmt.Errorf("query parsing error: %w", err)
		}

		return response, nil
	}
}

type StatPoint struct {
	Average float64 `json:"avg,omitempty"`
	Max     float64 `json:"max,omitempty"`
}

type cached struct {
	m          wrSVCMetrics
	logger     log.Logger
	getterFunc func() (map[string]map[string]map[time.Time]StatPoint, error)
	cachedData []byte
	mx         sync.Mutex
}

func (c *cached) GetStats(rw http.ResponseWriter, r *http.Request) {
	var err error
	defer func(begin time.Time) {
		lvs := []string{"method", "GetStats", "error", fmt.Sprint(err != nil)}
		c.m.requestCount.With(lvs...).Add(1)
		c.m.requestSize.With(lvs...).Add(float64(r.ContentLength))
		c.m.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	if c.cachedData == nil || r.Header.Get("Cache-Control") == "no-cache" {
		stats, err := c.getterFunc()
		if err != nil {
			http.Error(rw, err.Error(), 500)
			_ = c.logger.Log("svc", "wr", "during", "c.getterFunc()", "err", err)

			return
		}

		statsData, err := json.Marshal(&stats)
		if err != nil {
			http.Error(rw, err.Error(), 500)
			_ = c.logger.Log("svc", "wr", "during", "json.Marshal(&stats)", "err", err)

			return
		}

		c.mx.Lock()
		defer c.mx.Unlock()
		c.cachedData = statsData
	}
	rw.Header().Set("Cache-Control", "public, max-age=7200, immutable")
	_, _ = rw.Write(c.cachedData)
}

func CachedHTTPHandler(
	logger log.Logger,
	inflixDBConfig WriteAPIConfg,
) (*cached, func()) {
	client := influxdb2.NewClient(inflixDBConfig.url, inflixDBConfig.token)
	getterFunc := GetAllChampionStatsProvider(client, inflixDBConfig.bucket, inflixDBConfig.org)

	return &cached{
		m:          makeWRSVCMetrics(),
		logger:     logger,
		getterFunc: getterFunc,
	}, client.Close
}
