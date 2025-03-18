package client

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func parseKnockoutResults(res http.Response, compData CompetitionSeasonSummary) []MatchSummary {

	doc, err := goquery.NewDocumentFromReader(res.Body)
	
	if err != nil {
		log.Fatal("Failed to parse the HTML document", err)
	}
	
	goqueryString := "#sched_all tbody tr"

	var matchesData []MatchSummary

	doc.Find(goqueryString).Each(func(i int, row *goquery.Selection) {
		matchMap := MatchSummary{data: make(map[string]string), knockout: true}
		row.Find("th").Each(func(i int, cell *goquery.Selection) {
			statType, ok := cell.Attr("data-stat")
			if ok {
				processCell(statType, cell, &matchMap)
			}
		})
		row.Find("td").Each(func(i int, cell *goquery.Selection) {
			statType, ok := cell.Attr("data-stat")
			if ok {
				processCell(statType, cell, &matchMap)
			}
		})
		if (matchMap.data["home_team"] != "" || matchMap.data["away_team"] != "" || matchMap.data["score"] != "") &&  matchMap.data["score"] != "Score"  {
			matchesData = append(matchesData, matchMap)
		}
	},
	)
	return matchesData
}

func PrintKnockoutMatches(matches []MatchSummary, limit int) {
	for i, match := range matches {
		PrintKnockoutMatch(match, true)
		if i >= limit {
			return
		}
	}
}

func PrintKnockoutMatch(match MatchSummary, extended bool) {

	fmt.Println("===========================================================================================")
	fmt.Printf("%s match\n", match.data["round"])
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

