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

type GameIdentifier struct {
	RequestsUrl string
	Year        int
	GameInformation
}

type GameInformation struct {
	MatchDate                  string
	Stadium                    string
	Assitance                  string
	LocalTeam                  string
	LocalTeamYellowCards       int
	LocalTeamRedCards          int
	LocalTeamLeaguePosition    string
	LocalTeamballPosition      string
	LocalTeamScoreGoals        int
	VisitingTeam               string
	VisitingTeamYellowCards    int
	VisitingTeamRedCards       int
	VisitingTeamLeaguePosition string
	VisitingTeamballPosition   string
	VisitingTeamScoreGoals     int
	MatchResult                string
}

func Team_stats() {

	c := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))

	teamstats_collector := colly.NewCollector(colly.Async(true))
	teamstats_collector.Limit(&colly.LimitRule{Parallelism: 20})
	gameinformation_collector := teamstats_collector.Clone()

	Teams := []TeamStats{}
	Games := []GameIdentifier{}

	c.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("url", r.URL.String())
		fmt.Println("Visiting", r.URL.String())
	})

	teamstats_collector.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("url", r.URL.String())
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, "primera") && strings.Contains(link, "jornada") {
			visit_link := fmt.Sprintf("%s%s", helpers.Root_url, link)
			teamstats_collector.Visit(visit_link)
			teamstats_collector.Wait()
		}

		if strings.Contains(link, "partido") && !strings.Contains(link, "#videos") {
			visit_link := fmt.Sprintf("%s%s", helpers.Root_url, link)
			gameinformation_collector.Visit(visit_link)
			gameinformation_collector.Wait()
		}

	})

	teamstats_collector.OnHTML(`table`, func(e *colly.HTMLElement) {
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
						Teams = append(Teams, Team)
					}
					count += 1
				})
			})
		}
	})

	gameinformation_collector.OnHTML(`div[id=marcador]`, func(e *colly.HTMLElement) {
		request_url := e.Request.Ctx.Get("url")
		year := helpers.Parse_requests_url_year_game(request_url)
		GameIdentify := GameIdentifier{}
		GameIdentify.Year = year
		GameIdentify.RequestsUrl = request_url

		if GameIdentify.Year != 0 {
			Game := GameInformation{}
			Game.LocalTeamballPosition = e.ChildText(`div[class=header-team-1] 
											> div[class="team-stats posession"]`)
			Game.LocalTeamLeaguePosition = e.ChildText(`div[class=header-team-1] 
											> div[class="team-stats rank-pos"]`)
			Game.VisitingTeamballPosition = e.ChildText(`div[class=header-team-2] 
											> div[class="team-stats posession"]`)
			Game.VisitingTeamLeaguePosition = e.ChildText(`div[class=header-team-2] 
											> div[class="team-stats rank-pos"]`)
			Game.LocalTeam = e.ChildText(`div[class="team equipo1"] b`)
			Game.VisitingTeam = e.ChildText(`div[class="team equipo2"] b`)
			Game.MatchResult = e.ChildText(`div[class="resultado resultadoH"]`)
			Game.LocalTeamYellowCards = helpers.String_to_integer(
				e.ChildText(`div[class=te1] > span[class=am]`))
			Game.LocalTeamRedCards = helpers.String_to_integer(
				e.ChildText(`div[class=te1] > span[class=ro]`))
			Game.VisitingTeamYellowCards = helpers.String_to_integer(
				e.ChildText(`div[class=te2] > span[class=am]`))
			Game.VisitingTeamRedCards = helpers.String_to_integer(
				e.ChildText(`div[class=te2] > span[class=ro]`))
			Game.MatchDate = e.ChildText(`span[class=jor-date]`)
			Game.Stadium = e.ChildText(`div[class=matchinfo] li[class=es]`)
			Game.Assitance = e.ChildText(`div[class=matchinfo] li[class=as]`)
			Game.LocalTeamScoreGoals = helpers.String_to_integer(
				e.ChildText(`div[class="resultado resultadoH"]
												> span[class=claseR]:nth-child(1)`))
			Game.VisitingTeamScoreGoals = helpers.String_to_integer(
				e.ChildText(`div[class="resultado resultadoH"]
								> span[class=claseR]:nth-child(2)`))
			GameIdentify.GameInformation = Game
			Games = append(Games, GameIdentify)
		}
	})

	for year := 1932; year <= 2019; year++ {
		root_url := helpers.Scrapper_root_url
		start_date := fmt.Sprintf("%s/%s", strconv.Itoa(year), "/grupo1/jornada")
		visit_link := fmt.Sprintf("%s%s", root_url, start_date)
		c.Visit(visit_link)

	}

	TeamsWriter, _ := os.Create("scrapper_results/results_evolution.json")
	json.NewEncoder(TeamsWriter).Encode(Teams)

	GamesWriter, _ := os.Create("scrapper_results/games_historical.json")
	json.NewEncoder(GamesWriter).Encode(Games)
}
