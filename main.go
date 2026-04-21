package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// ... (Estructuras OddsResponse, Bookmaker, Market, Outcome se mantienen igual) ...
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
	Key      string    `json:"key"`
	Outcomes []Outcome `json:"outcomes"`
}
type Outcome struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Point float64 `json:"point"`
}

func main() {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		fmt.Println("Error: API_KEY no configurada")
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
	sb.WriteString(".game { background: #1e1e1e; padding: 15px; margin-bottom: 15px; border-radius: 8px; border-left: 5px solid #00ff88; }")
	sb.WriteString(".bookie-name { color: #00ff88; font-weight: bold; margin-top: 5px; display: block; }")
	sb.WriteString(".odds-row { font-size: 0.85em; color: #ccc; margin-left: 10px; padding: 2px 0; }")
	sb.WriteString(".val { color: #fff; font-weight: bold; }")
	sb.WriteString("h1, h2 { color: #00ff88; text-transform: uppercase; border-bottom: 1px solid #333; padding-bottom: 5px; }")
	sb.WriteString("</style></head><body>")
	sb.WriteString(fmt.Sprintf("<h1>📅 %s (UTC)</h1>", time.Now().Format("Jan 02, 15:04")))

	for _, sport := range sports {
		// Pedimos ambos mercados para todos, pero luego filtraremos en el HTML
		url := fmt.Sprintf("https://api.the-odds-api.com/v4/sports/%s/odds/?apiKey=%s&regions=us&markets=h2h,spreads&oddsFormat=decimal", sport, apiKey)

		resp, err := http.Get(url)
		if err != nil {
			continue
		}

		var games []OddsResponse
		json.NewDecoder(resp.Body).Decode(&games)
		resp.Body.Close()

		sb.WriteString(fmt.Sprintf("<h2>%s %s</h2>", emojis[sport], strings.Replace(sport, "_", " ", -1)))

		for _, game := range games {
			sb.WriteString("<div class='game'>")
			sb.WriteString(fmt.Sprintf("<strong>%s vs %s</strong>", game.AwayTeam, game.HomeTeam))

			for _, bm := range game.Bookmakers {
				sb.WriteString(fmt.Sprintf("<span class='bookie-name'>%s:</span>", bm.Title))

				for _, market := range bm.Markets {

					// --- LÓGICA DE FILTRADO ---
					// Si es NBA, mostramos Spread. Si NO es NBA, mostramos solo H2H (Ganador).
					if sport == "basketball_nba" && market.Key != "spreads" {
						continue // En NBA saltamos el H2H
					}
					if sport != "basketball_nba" && market.Key != "h2h" {
						continue // En otros deportes saltamos el Spread
					}
					// ---------------------------

					sb.WriteString("<div class='odds-row'>")
					label := "🏆 Ganador: "
					if market.Key == "spreads" {
						label = "⚖️ Spread: "
					}
					sb.WriteString(label)

					for i, outcome := range market.Outcomes {
						pointStr := ""
						if market.Key == "spreads" {
							pointStr = fmt.Sprintf("(%+.1f) ", outcome.Point)
						}
						sb.WriteString(fmt.Sprintf("%s %s<span class='val'>%.2f</span>", outcome.Name, pointStr, outcome.Price))
						if i < len(market.Outcomes)-1 {
							sb.WriteString(" | ")
						}
					}
					sb.WriteString("</div>")
				}
			}
			sb.WriteString("</div>")
		}
		time.Sleep(500 * time.Millisecond)
	}

	sb.WriteString("</body></html>")
	os.WriteFile("index.html", []byte(sb.String()), 0644)
}
