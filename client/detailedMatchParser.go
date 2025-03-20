package client

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func processDocAndTeam(res *http.Response, homePlayers *[]PlayerDetailContainer, awayPlayers *[]PlayerDetailContainer, homeCode string, awayCode string) {

	// defer res.Body.Close()
	
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("Failed to parse the HTML document", err.Error())
		os.Exit(1)
	}

	homeSearch := fmt.Sprintf("#stats_%s_summary tbody tr ", homeCode)
	awaySearch := fmt.Sprintf("#stats_%s_summary tbody tr ", awayCode)

	// homeSearch := "#stats_bba7d733_summary tbody tr"
	// awaySearch := "#stats_6ca73159_summary tbody tr"

	doc.Find(homeSearch).Each(func(i int, row *goquery.Selection) {

		player := PlayerDetailContainer{involved_in_substitution: false}

		//this will select and process the player name
		row.Find("th").Each(func(j int, cell *goquery.Selection) {
			statType, ok := cell.Attr("data-stat")
			if ok {
				processPlayerCell(statType, cell, &player)
			}
		})
		//this will select and process all other player attributes
		row.Find("td").Each(func(j int, cell *goquery.Selection) {
			statType, ok := cell.Attr("data-stat")
			if ok {
				processPlayerCell(statType, cell, &player)
			}
		})

		*homePlayers = append(*homePlayers, player)
	
	})
	doc.Find(awaySearch).Each(func(i int, row *goquery.Selection) {

		player := PlayerDetailContainer{involved_in_substitution: false}

		//this will select and process the player name
		row.Find("th").Each(func(j int, cell *goquery.Selection) {
			statType, ok := cell.Attr("data-stat")
			if ok {
				processPlayerCell(statType, cell, &player)
			}
		})
		//this will select and process all other player attributes
		row.Find("td").Each(func(j int, cell *goquery.Selection) {
			statType, ok := cell.Attr("data-stat")
			if ok {
				processPlayerCell(statType, cell, &player)
			}
		})

		*awayPlayers = append(*awayPlayers, player)

	})}

func processPlayerCell(stat string, cell *goquery.Selection, player *PlayerDetailContainer) {

	switch stat {
	case "player":
		player.player_name = cell.Find("a").Text()
		link, ok := cell.Find("a").Attr("href")
		if ok {
			player.player_url = link
		}
	case "nationality":
		player.nationality = cell.Find("a").Text()
		n, ok := cell.Find("a").Attr("href")
		if ok {
			player.nationality_link = n
		}
	case "age":
		age := cell.Text()
		player.ageYears = age[0:2]
		player.ageDays = age[3:]
	case "minutes":
		if mins := cell.Text(); mins == "90" {
			player.mins_played = 90
			player.involved_in_substitution = false
		} else {
			player.involved_in_substitution = true
			nMins, err := strconv.Atoi(mins)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			player.mins_played = nMins
		}
	case "goals":
		gls, gErr := strconv.Atoi(cell.Text())
		if gErr != nil {
			fmt.Println(gErr.Error())
			os.Exit(1)
		}
		player.goals = gls
	case "pens_made":
		pens, pErr := strconv.Atoi(cell.Text())
		if pErr != nil {
			fmt.Println(pErr.Error())
			os.Exit(1)
		}
		player.penalties = pens
	case "cards_yellow":
		ycs, yErr := strconv.Atoi(cell.Text())
		if yErr != nil {
			fmt.Println(yErr.Error())
			os.Exit(1)
		}
		player.yellow_card = ycs
	case "cards_red":
		rcs, rErr := strconv.Atoi(cell.Text())
		if rErr != nil {
			fmt.Println(rErr.Error())
			os.Exit(1)
		}
		player.red_card = rcs
	case "own_goals":
		ogs, oErr := strconv.Atoi(cell.Text())
		if oErr != nil {
			fmt.Println(oErr.Error())
			os.Exit(1)
		}
		player.own_goals = ogs
	default:
		return
	}
}

type DetailedMatchContainer struct {
	matchUrl, HomeTeamOnlineID, AwayTeamOnlineID string
	homePlayers, awayPlayers []PlayerDetailContainer
	isknockout bool
}

type PlayerDetailContainer struct {
	//player identifiers
	player_name string
    nationality string
	nationality_link string
    ageYears string
	ageDays string
    player_url string
	//match statistics
	first_minute int
	last_minute int
	mins_played int
	involved_in_substitution bool
    goals int
    penalties int
    yellow_card int
    red_card int
    own_goals int
	home_or_away bool
}