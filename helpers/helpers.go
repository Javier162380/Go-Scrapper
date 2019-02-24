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

func Parse_request_url_player_year(requests_input string) int {
	if strings.Contains(requests_input, "jugador/") {
		root_url := strings.Split(requests_input, "jugador/")
		if len(root_url) > 1 {
			root_url_split := root_url[1]
			year_split := strings.Split(root_url_split, "/")
			if len(year_split) > 1 {
				return String_to_integer(year_split[0])
			}
		}
	}
	return 0
}

func Parse_request_url_player_name(requests_input string) string {
	if strings.Contains(requests_input, "jugador/") {
		root_url := strings.Split(requests_input, "jugador/")
		if len(root_url) > 1 {
			root_url_split := root_url[1]
			year_split := strings.Split(root_url_split, "/")
			if len(year_split) > 1 {
				player_split := strings.Split(root_url_split, "-")
				if len(player_split) > 0 {
					player_year_split := strings.Split(player_split[0], "/")
					if len(player_year_split) > 1 {
						return player_year_split[1]
					}
					return player_split[0]
				}
			}
		}
	}
	return ""

}

func Parse_request_url_player_id(requests_input string) int {
	if strings.Contains(requests_input, "jugador/") {
		root_url := strings.Split(requests_input, "jugador/")
		if len(root_url) > 1 {
			root_url_split := root_url[1]
			year_split := strings.Split(root_url_split, "/")
			if len(year_split) > 1 {
				player_split := strings.Split(root_url_split, "-")
				if len(player_split) > 1 {
					return String_to_integer(player_split[len(player_split)-1])
				}
			}
		}
	}
	return 0

}
