package client

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func processDocAndTeam(res *http.Response, homePlayers *[]PlayerDetailContainer, awayPlayers *[]PlayerDetailContainer, homeCode string, awayCode string, rowUrl string) {

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
				err := processPlayerCell(statType, cell, &player)
				if err != nil {
					return
				}
			}
		})
		//this will select and process all other player attributes
		row.Find("td").Each(func(j int, cell *goquery.Selection) {
			statType, ok := cell.Attr("data-stat")
			if ok {
				err := processPlayerCell(statType, cell, &player)
				if err != nil {
					return
				}
			}
		})

		*awayPlayers = append(*awayPlayers, player)

	})}

func processPlayerCell(stat string, cell *goquery.Selection, player *PlayerDetailContainer) error {

	switch stat {
	case "player":
		player.player_name = cell.Find("a").Text()
		link, ok := cell.Find("a").Attr("href")
		if ok {
			player.player_url = link
		} else {
			return fmt.Errorf("issue with accessing the player name / url")
		}
	case "nationality":
		player.nationality = cell.Find("a").Text()
		n, ok := cell.Find("a").Attr("href")
		if ok {
			player.nationality_link = n
		} else {
			return fmt.Errorf("issue with accessing the player nationality")
		}
	case "age":
		age := cell.Text()
		if len(age) == 6 {
			player.ageYears = age[0:2]
			player.ageDays = age[3:]
		} else {return fmt.Errorf("couldn't access player age")}
	case "minutes":
		if mins := cell.Text(); mins == "90" {
			player.mins_played = 90
			player.involved_in_substitution = false
		} else {
			player.involved_in_substitution = true
			nMins, err := strconv.Atoi(mins)
			if err != nil {
				return err
			}
			player.mins_played = nMins
		}
	case "goals":
		gls, gErr := strconv.Atoi(cell.Text())
		if gErr != nil {
			return gErr
		}
		player.goals = gls
	case "pens_made":
		pens, pErr := strconv.Atoi(cell.Text())
		if pErr != nil {
			return pErr
		}
		player.penalties = pens
	case "cards_yellow":
		ycs, yErr := strconv.Atoi(cell.Text())
		if yErr != nil {
			return yErr
		}
		player.yellow_card = ycs
	case "cards_red":
		rcs, rErr := strconv.Atoi(cell.Text())
		if rErr != nil {
			return rErr
		}
		player.red_card = rcs
	case "own_goals":
		ogs, oErr := strconv.Atoi(cell.Text())
		if oErr != nil {
			return oErr
		}
		player.own_goals = ogs
	default:
		return nil
	}
	return nil
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