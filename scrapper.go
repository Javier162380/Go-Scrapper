package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

type MatchStats struct {
	Date         string
	LocalTeam    string
	Result       string
	VisitingTeam string
	Stadium      string
}

type Clasification struct {
	Position int
	TeamStats
}

type TeamStats struct {
	Position      int
	Team          string
	Points        int
	League_Day    int
	MatchesPlayed int
	MatchesWin    int
	MatchesDraw   int
	MatchesLoose  int
	GoalsScore    int
	GoalsRecieve  int
}

func string_to_integer(scraper_input string) int {
	integer, err := strconv.Atoi(scraper_input)
	if err != nil {
	}
	return integer
}
func main() {
	c := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c.OnHTML(`table`, func(e *colly.HTMLElement) {
		table_id := e.Attr(`id`)
		if strings.Contains(table_id, "tabla1") {
			Matches_list := []MatchStats{}
			e.ForEach("table tbody", func(_ int, el *colly.HTMLElement) {
				ch := e.DOM.Children()
				Match := MatchStats{}
				ch.Find("tr").Each(func(clss int, tr *goquery.Selection) {
					row_node := tr.Find("td")
					Match.Result = row_node.Find("span").Last().Text()
					Match.Date = row_node.Find("span").Next().First().Text()
					Match.Stadium = row_node.Find("span").Next().Next().First().Text()
					Match.LocalTeam = row_node.Find("a").Next().First().Text()
					Match.VisitingTeam = row_node.Find("a").Next().Slice(1, 2).Text()
					Matches_list = append(Matches_list, Match)
				})
			})
			fmt.Println("%s", Matches_list)

		}

		if strings.Contains(table_id, "tabla2") {
			Table := []TeamStats{}
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
						Team.Position = string_to_integer(index_node.Text())
						Team.Team = row_node.Find("a").First().Text()
						Team.Points = string_to_integer(row_node.Next().First().Text())
						Team.MatchesPlayed = string_to_integer(row_node.Next().Next().First().Text())
						Team.MatchesWin = string_to_integer(row_node.Next().Next().Next().First().Text())
						Team.MatchesDraw = string_to_integer(row_node.Next().Next().Next().Next().First().Text())
						Team.MatchesLoose = string_to_integer(row_node.Next().Next().Next().Next().Next().First().Text())
						Team.GoalsScore = string_to_integer(row_node.Next().Next().Next().Next().Next().Next().Text())
						Team.GoalsRecieve = string_to_integer(row_node.Next().Next().Next().Next().Next().Next().Next().First().Text())
						Table = append(Table, Team)
					}
					count += 1
				})
			})
			fmt.Println("%s", Table)
		}
	})
	c.Visit("http://www.resultados-futbol.com/primera1932/grupo1/jornada1")

}
