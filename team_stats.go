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

type TeamStats struct {
	Year            int
	LeagueDate      int
	Position        int
	Team            string
	Points          int
	MatchesPlayed   int
	MatchesWin      int
	MatchesDraw     int
	MatchesLoose    int
	GoalsScore      int
	GoalsRecieve    int
	GoalsDifference int
}

func Teamstats() {

	c := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))

	Table := []TeamStats{}

	c.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("url", r.URL.String())
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, "primera") && strings.Contains(link, "jornada") {
			e.Request.Visit(link)
		}

	})

	c.OnHTML(`table`, func(e *colly.HTMLElement) {
		table_id := e.Attr(`id`)
		request_url := e.Request.Ctx.Get("url")
		LeagueDate_Split := strings.Split(request_url, "/grupo1/jornada")
		Year_Split := strings.Split(LeagueDate_Split[0], "primera")[1]
		Year := helpers.Parse_request_url_year(Year_Split)
		if strings.Contains(table_id, "tabla2") {
			e.ForEach("table tbody", func(_ int, el *colly.HTMLElement) {
				ch := e.DOM.Children()
				Team := TeamStats{}
				count := 0
				ch.Find("tr").Each(func(clss int, tr *goquery.Selection) {
					if count == 0 {
					}
					if count != 0 {
						row_node := tr.Find("td").First()
						index_node := tr.Find("th")
						Team.Year = helpers.String_to_integer(Year)
						Team.LeagueDate = helpers.String_to_integer(row_node.Next().Next().First().Text())
						Team.Position = helpers.String_to_integer(index_node.Text())
						Team.Team = row_node.Find("a").First().Text()
						Team.Points = helpers.String_to_integer(row_node.Next().First().Text())
						Team.MatchesPlayed = helpers.String_to_integer(row_node.Next().Next().First().Text())
						Team.MatchesWin = helpers.String_to_integer(row_node.Next().Next().Next().First().Text())
						Team.MatchesDraw = helpers.String_to_integer(row_node.Next().Next().Next().Next().First().Text())
						Team.MatchesLoose = helpers.String_to_integer(row_node.Next().Next().Next().Next().Next().First().Text())
						Team.GoalsScore = helpers.String_to_integer(row_node.Next().Next().Next().Next().Next().Next().Text())
						Team.GoalsRecieve = helpers.String_to_integer(row_node.Next().Next().Next().Next().Next().Next().Next().First().Text())
						Team.GoalsDifference = Team.GoalsScore - Team.GoalsRecieve
						Table = append(Table, Team)
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

	resultsWriter, _ := os.Create("scrapper_results/results_evolution_v3.json")
	json.NewEncoder(resultsWriter).Encode(Table)

}
