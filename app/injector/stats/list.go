package stats

import (
	"fmt"
	"net/http"

	"github.com/bibaroc/paladinswr/app/injector/heroes"
)

func ListStats() ([]ShortChampionStats, error) {
	statsClient := PaladinsGuru(http.DefaultClient)
	allHeroes := heroes.All()
	shortStats := make([]ShortChampionStats, len(allHeroes))

	for i := range allHeroes {
		maxLoadoutWinrate := 0.0

		totalPlayed := 0
		totalWon := 0

		heroStat, err := statsClient.GetChampionStats(allHeroes[i])
		if err != nil {
			return nil, fmt.Errorf("failed to get stats for %s: %w", allHeroes[i], err)
		}

		for ii := range heroStat.Loadouts {
			loadoutWinrate := float64(heroStat.Loadouts[ii].Wins) / float64(heroStat.Loadouts[ii].Played) * 100
			if loadoutWinrate > maxLoadoutWinrate {
				maxLoadoutWinrate = loadoutWinrate
			}

			totalPlayed += heroStat.Loadouts[ii].Played
			totalWon += heroStat.Loadouts[ii].Wins
		}

		shortStats[i].Name = heroStat.Champion.Name
		shortStats[i].Class = heroStat.Champion.Class
		shortStats[i].Winrate.BestLoadout = maxLoadoutWinrate
		shortStats[i].Winrate.WeightedAverage = float64(totalWon) / float64(totalPlayed) * 100

	}

	return shortStats, nil
}
