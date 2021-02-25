package stats

import "time"

// ChampionStats is generated via `https://mholt.github.io/json-to-go/`
// based on output from
// GET	https://api.paladins.guru/v3/champions/zhin/builds
type ChampionStats struct {
	Champion struct {
		Abilities struct {
			Num1 struct {
				Name string `json:"name"`
			} `json:"1"`
			Num2 struct {
				Name string `json:"name"`
			} `json:"2"`
			Num3 struct {
				Name string `json:"name"`
			} `json:"3"`
			Num4 struct {
				Name string `json:"name"`
			} `json:"4"`
			Num5 struct {
				Name string `json:"name"`
			} `json:"5"`
		} `json:"abilities"`
		Name  string `json:"name"`
		Title string `json:"title"`
		Class string `json:"class"`
		ID    int    `json:"id"`
	} `json:"champion"`
	Meta struct {
		Champions string `json:"champions"`
		Items     string `json:"items"`
		Patch     struct {
			Start    time.Time `json:"start"`
			ID       string    `json:"_id"`
			Version  string    `json:"version"`
			Platform int       `json:"platform"`
		} `json:"patch"`
	} `json:"_meta"`
	Loadouts []struct {
		ID    string `json:"_id"`
		Cards []struct {
			ID    int `json:"id"`
			Level int `json:"level"`
		} `json:"cards"`
		Played int `json:"played"`
		Wins   int `json:"wins"`
	} `json:"loadouts"`
	Builds struct {
		Actives []struct {
			Wl     float64 `json:"wl"`
			Wins   int     `json:"wins"`
			Item   int     `json:"item"`
			Usage  float64 `json:"usage"`
			Played int     `json:"played"`
		} `json:"actives"`
		Talents []struct {
			Wl     float64 `json:"wl"`
			Wins   int     `json:"wins"`
			Item   int     `json:"item"`
			Usage  float64 `json:"usage"`
			Played int     `json:"played"`
		} `json:"talents"`
		ProBuilds []struct {
			Player struct {
				Name     string `json:"name"`
				League   string `json:"league"`
				Team     string `json:"team"`
				PlayerID int    `json:"player_id"`
			} `json:"player"`
			Loadout []struct {
				ID    int `json:"id"`
				Level int `json:"level"`
			} `json:"loadout"`
			Actives []struct {
				ID    int `json:"id"`
				Level int `json:"level"`
			} `json:"actives"`
			Talent int `json:"talent"`
		} `json:"pro_builds"`
		Played int `json:"played"`
	} `json:"builds"`
}
