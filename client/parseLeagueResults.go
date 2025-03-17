package client

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseLeagueResults(res http.Response, compData CompetitionSeasonSummary) []MatchSummary {

	doc, err := goquery.NewDocumentFromReader(res.Body)
		
	if err != nil {
		log.Fatal("Failed to parse the HTML document", err)
	}
	
	goqueryString := fmt.Sprintf("#sched_%s_%s_1 tbody tr",
		compData.CompetitionSeason,
		compData.CompetitionOnlineID,
	)

	var matchesData []MatchSummary

	doc.Find(goqueryString).Each(func(i int, row *goquery.Selection) {
		matchMap := MatchSummary{data: make(map[string]string)}
		row.Find("td").Each(func(i int, cell *goquery.Selection) {
			statType, ok := cell.Attr("data-stat")
			if ok {
				processCell(statType, cell, &matchMap)
			}
		})
		matchesData = append(matchesData, matchMap)
	},
	)
	return matchesData
}

func processCell(statType string, cell *goquery.Selection, summary *MatchSummary) {

	switch statType {
	case "home_team", "away_team", "date":
		summary.data[statType] = cell.Text()
		link, ok := cell.Find("a").Attr("href")
		if ok {
			summary.data[fmt.Sprintf("%s-url",statType)] = link
		}
	case "score":
		summary.data[statType] = cell.Text()
		link, ok := cell.Find("a").Attr("href")
		if ok {
			summary.data["match-report"] = link
		}
	case "notes", "match_report":
		return
	default:
		summary.data[statType] = cell.Text()
	}}


func PrintMatches(matches []MatchSummary, limit int) {
	for i, match := range matches {
		PrintMatch(match, true)
		if i >= limit {
			return
		}
	}
}

func PrintMatch(match MatchSummary, extended bool) {

	fmt.Println("===========================================================================================")
	fmt.Printf("%s %s %s\n", match.data["home_team"], match.data["score"], match.data["away_team"])
	fmt.Printf("Time: %s, %s, %s\n", match.data["start_time"], match.data["dayofweek"], match.data["date"])
	fmt.Printf("Referee: %s\n", match.data["referee"])
	fmt.Printf("Venue: %s, Attendance: %s\n", match.data["venue"], match.data["attendance"])
	fmt.Println("===========================================================================================")


	if extended {
		for key, val := range match.data {
			key = strings.Repeat(" ", 20-len(key)) + key 
			fmt.Printf("%s : %s\n",key, val)
		}
	}
	fmt.Println()
}

type MatchSummary struct {
	data map[string]string
}

type CompetitionSeasonSummary struct {
	Data []MatchSummary
	CompetitionName string
	CompetitionSeason string
	CompetitionOnlineID string
	Url string
}