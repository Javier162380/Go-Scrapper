package soccer_scrapper

import (
	"encoding/json"
	"fmt"
	"javier162380/soccer-scrapper/helpers"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

type Player struct {
	PlayerID int
	Name     string
	Year     int
	PlayerStats
}

type PlayerStats struct {
	Competition     string
	Team            string
	GamesPlayed     int
	GamesStarting   int
	GamesCompleted  int
	GamesSustituted int
	MinutesPlayed   string
	YellowCards     int
	RedCards        int
	Assits          int
	GoalsScore      int
}

func main() {
	c := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))

	url_collector := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}),
		colly.Async(true))
	players_collector := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}),
		colly.Async(true))
	Players_List := []Player{}

	players_collector.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("url", r.URL.String())
	})

	url_collector.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("url", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, "jugador/") {
			visit_link := fmt.Sprintf("%s%s", helpers.Root_url, link)
			url_collector.Visit(visit_link)
			url_collector.Wait()
			url_collector.Limit(&colly.LimitRule{Parallelism: 20})
		}

	})

	url_collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		requests_url := e.Request.Ctx.Get("url")
		year := helpers.Parse_request_url_player_year(requests_url)
		if strings.Contains(link, fmt.Sprintf("%s/%s", "jugador", strconv.Itoa(year))) {
			players_collector.Visit(link)
			players_collector.Wait()
			players_collector.Limit(&colly.LimitRule{
				Parallelism: 20})

		}
	})

	players_collector.OnHTML(`table`, func(e *colly.HTMLElement) {
		table_id := e.Attr(`class`)
		requests_url := e.Request.Ctx.Get("url")
		if strings.Contains(table_id, "h-classification") {
			e.ForEach("table tbody", func(_ int, el *colly.HTMLElement) {
				Player_Struct := Player{}
				year := helpers.Parse_request_url_player_year(requests_url)
				name := helpers.Parse_request_url_player_name(requests_url)
				id := helpers.Parse_request_url_player_id(requests_url)
				Player_Struct.Name = name
				Player_Struct.Year = year
				Player_Struct.PlayerID = id
				ch := e.DOM.Children()
				count := 0
				ch.Find("tr").Each(func(clss int, tr *goquery.Selection) {
					if count == 0 {
					}
					if count != 0 {
						Stats := PlayerStats{}
						headers := tr.Find("th")
						row_node := tr.Find("td")
						Competition := headers.Find("span").Next().First().Text()
						if strings.Contains(Competition, "Primera División") {
							Stats.Competition = headers.Find("span").Next().First().Text()
							Stats.Team = strings.TrimSpace(headers.Siblings().First().Text())
							Stats.GamesPlayed = helpers.String_to_integer(row_node.First().Text())
							Stats.GamesStarting = helpers.String_to_integer(row_node.Next().First().Text())
							Stats.GamesCompleted = helpers.String_to_integer(row_node.Next().Next().First().Text())
							Stats.GamesSustituted = helpers.String_to_integer(row_node.Next().Next().Next().First().Text())
							Stats.MinutesPlayed = row_node.Next().Next().Next().Next().First().Text()
							Stats.YellowCards = helpers.String_to_integer(row_node.Next().Next().Next().Next().Next().First().Text())
							Stats.RedCards = helpers.String_to_integer(row_node.Next().Next().Next().Next().Next().Next().First().Text())
							Stats.Assits = helpers.String_to_integer(row_node.Next().Next().Next().Next().Next().Next().Next().First().Text())
							Stats.GoalsScore = helpers.String_to_integer(row_node.Next().Next().Next().Next().Next().Next().Next().Next().First().Text())
							Player_Struct.PlayerStats = Stats
							Players_List = append(Players_List, Player_Struct)
						}

					}
					count += 1
				})
			})
		}
	})

	for year := 1932; year <= 2019; year++ {
		root_url := helpers.Scrapper_root_url
		start_date := fmt.Sprintf("%s/%s", strconv.Itoa(year), "/grupo1/jornada")
		visit_link := fmt.Sprintf("%s%s", root_url, start_date)
		c.Visit(visit_link)
	}

	resultsWriter, _ := os.Create("scrapper_results/players_json.json")
	json.NewEncoder(resultsWriter).Encode(Players_List)

}
