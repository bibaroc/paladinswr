package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bibaroc/paladinswr/heroes"
	"github.com/bibaroc/paladinswr/stats"
)

// influxdb2 "github.com/influxdata/influxdb-client-go/v2"

func main() {
	statsClient := stats.PaladinsGuru(http.DefaultClient)

	for _, name := range heroes.Flanker() {
		heroStat, err := statsClient.GetChampionStats(name)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(heroStat.Champion.Name)

		maxLoadoutWinrate := 0.0

		for i := range heroStat.Loadouts {
			loadoutWinrate := float64(heroStat.Loadouts[i].Wins) / float64(heroStat.Loadouts[i].Played) * 100
			if loadoutWinrate > maxLoadoutWinrate {
				maxLoadoutWinrate = loadoutWinrate
			}
		}

		fmt.Println("\t", "maxLoadoutWinrate", maxLoadoutWinrate)
	}
}
