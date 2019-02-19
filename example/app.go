package main

import (
	"fmt"
	"javier162380/soccer-scrapper/helpers"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

func main() {
	c := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))

	c.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("url", r.URL.String())
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Println(link)
		//request_url := e.Request.Ctx.Get("url")
		//LeagueDate_Split := strings.Split(request_url, "/grupo1/jornada")
		//Year_Split := strings.Split(LeagueDate_Split[0], "primera")[1]
		//Root_Url_Year := helpers.String_to_integer(helpers.Parse_request_url_year(Year_Split))
		if strings.Contains(link, "jugador/") {
			link_year := helpers.Parse_request_url_player(link)
			if link_year == 1932 {
				//fmt.Println(link)
				//e.Request.Visit(link)
			}
		}

	})

	c.OnHTML(`table`, func(e *colly.HTMLElement) {
		table_id := e.Attr(`class`)
		if strings.Contains(table_id, "h-classification") {
			fmt.Println("aqui estamos")
			e.ForEach("table tbody", func(_ int, el *colly.HTMLElement) {
				ch := e.DOM.Children()
				ch.Find("tr").Each(func(clss int, tr *goquery.Selection) {
				})
			})
		}
	})

	c.Visit("http://www.resultados-futbol.com/primera1932/grupo1/jugadores")

}
