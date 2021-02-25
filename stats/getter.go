package stats

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const buildsURLTemplate = "https://api.paladins.guru/v3/champions/%s/builds"

type paladinsGuruStat struct {
	client *http.Client
}

func (pgs *paladinsGuruStat) GetChampionStats(champion string) (ChampionStats, error) {
	getBuildResponse, err := pgs.client.Get(fmt.Sprintf(buildsURLTemplate, url.PathEscape(champion)))
	if err != nil {
		return ChampionStats{}, fmt.Errorf("couldn't query builds for %s: %w", champion, err)
	}

	defer getBuildResponse.Body.Close()

	stats := ChampionStats{}
	if err = json.NewDecoder(getBuildResponse.Body).Decode(&stats); err != nil {
		return ChampionStats{}, fmt.Errorf("failed to unmarshal stats: %w", err)
	}

	return stats, nil
}

func PaladinsGuru(client *http.Client) *paladinsGuruStat {
	return &paladinsGuruStat{client: client}
}
