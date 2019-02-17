package helpers

import (
	"strconv"
	"strings"
)

func String_to_integer(scraper_input string) int {
	integer, err := strconv.Atoi(scraper_input)
	if err != nil {
	}
	return integer
}

func Parse_request_url_year(requests_input string) string {
	if strings.Contains(requests_input, "/") {
		return strings.Split(requests_input, "/")[0]
	}
	return requests_input
}
