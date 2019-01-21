package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

type Matches struct {
	Date         string
	LocalTeam    string
	Result       string
	VisitingTeam string
	Location     string
}

type Clasification struct {
	Position int
	TeamStats
}

type TeamStats struct {
	Team          string
	Points        int
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
			Matches_list := []Matches{}
			e.ForEach("table tbody", func(_ int, el *colly.HTMLElement) {
				ch := e.DOM.Children()
				Match := Matches{}
				ch.Find("tr").Each(func(clss int, tr *goquery.Selection) {
					row_node := tr.Find("td")
					if row_node.HasClass("equipo1") {
						row_node.Find(".equipo1").Each(func(clss int, s *goquery.Selection) {
							band := s.Attr("class")
							fmt.Printf("%s", band)
						})

						Match.LocalTeam = row_node.Find(".equipo1").Text()

						fmt.Printf("%s", "Hola")
						//Match.LocalTeam = row_node.Find("a:nth-child(2)").Text()
						//clase := row_node.RemoveClass("equipo1")
						//fmt.Printf("%s", row_node.Find("a:nth-of-type(3)").Text())
						//Find("a:nth-child(2)").Text())
					}
					//if row_node.HasClass("equipo2") {

					//	Match.VisitingTeam = row_node.Find("a:nth-child(2)").Text()
					//}
					//if row_node.HasClass("rstd") {
					//	row_node.Find(".rstd").Each(func(_ int, ul *goquery.Selection) {
					//		a := row_node.Find(".rstd")
					//		fmt.Printf("%s", a.Text())
					//	})
					//Match.Result = clss
					//	Match.Date = row_node.Find("span:nth-of-type(2)").Text()
					//	Match.Location = row_node.Find("span:nth-of-type(3)").Text()
					//fmt.Printf("%s", Match.Result)
					//}
					//fmt.Printf("%s", Match)
					Matches_list = append(Matches_list, Match)
				})
			})
		}

		if strings.Contains(table_id, "tabla2") {
			LeagueTable := Clasification{}
			e.ForEach("table tbody", func(_ int, el *colly.HTMLElement) {
				ch := e.DOM.Children()
				ch.Find("tr").Each(func(clss int, tr *goquery.Selection) {
					LeagueTable.Position = string_to_integer(tr.Find("th").Text())
					time.Sleep(10)
					row_node := tr.Find("td")
					if row_node.HasClass("equipo") {
						LeagueTable.Team = row_node.Find("a").Text()
					}
					//fmt.Printf("%s", LeagueTable.Position)
					//if row_node.HasClass("pts") {
					//	LeagueTable.Points = row_node.Attr("")
					//}

				})
			})

		}

	})

	c.Visit("http://www.resultados-futbol.com/primera1932/grupo1/jornada1")

}
