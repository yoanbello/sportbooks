package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type OddsResponse struct {
	HomeTeam   string      `json:"home_team"`
	AwayTeam   string      `json:"away_team"`
	Bookmakers []Bookmaker `json:"bookmakers"`
}

type Bookmaker struct {
	Title   string   `json:"title"`
	Markets []Market `json:"markets"`
}

type Market struct {
	Outcomes []Outcome `json:"outcomes"`
}

type Outcome struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func main() {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		fmt.Println("Error: No se encontró la API_KEY")
		return
	}

	sports := []string{"basketball_nba", "icehockey_nhl", "baseball_mlb"}
	emojis := map[string]string{
		"basketball_nba": "🏀",
		"icehockey_nhl":  "🏒",
		"baseball_mlb":   "⚾",
	}

	var sb strings.Builder
	sb.WriteString("<html><head><meta charset='UTF-8'><style>")
	sb.WriteString("body { font-family: sans-serif; background: #121212; color: #e0e0e0; padding: 20px; }")
	sb.WriteString(".game { background: #1e1e1e; padding: 15px; margin-bottom: 10px; border-radius: 8px; border-left: 5px solid #00ff88; }")
	sb.WriteString(".bookie { font-size: 0.9em; color: #aaa; margin-left: 15px; padding: 2px; }")
	sb.WriteString(".odds { font-weight: bold; color: #00ff88; }")
	sb.WriteString("h1 { color: #00ff88; }")
	sb.WriteString("</style></head><body>")
	sb.WriteString(fmt.Sprintf("<h1>Apuestas Actualizadas: %s</h1>", time.Now().Format("15:04:05")))

	for _, sport := range sports {
		url := fmt.Sprintf("https://api.the-odds-api.com/v4/sports/%s/odds/?apiKey=%s&regions=us&markets=h2h", sport, apiKey)
		resp, err := http.Get(url)
		if err != nil {
			continue
		}

		var games []OddsResponse
		json.NewDecoder(resp.Body).Decode(&games)
		resp.Body.Close()

		emoji := emojis[sport]
		sb.WriteString(fmt.Sprintf("<h2>%s %s</h2>", emoji, strings.ToUpper(sport)))

		for _, game := range games {
			sb.WriteString("<div class='game'>")
			sb.WriteString(fmt.Sprintf("<strong>%s vs %s</strong><br>", game.AwayTeam, game.HomeTeam))
			for _, bm := range game.Bookmakers {
				sb.WriteString(fmt.Sprintf("<div class='bookie'>[%s] ", bm.Title))
				for _, outcome := range bm.Markets[0].Outcomes {
					sb.WriteString(fmt.Sprintf("%s: <span class='odds'>%.2f</span> | ", outcome.Name, outcome.Price))
				}
				sb.WriteString("</div>")
			}
			sb.WriteString("</div>")
		}
		time.Sleep(500 * time.Millisecond)
	}
	sb.WriteString("</body></html>")
	os.WriteFile("index.html", []byte(sb.String()), 0644)
}
