package client

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func ParseLeagueResults(res http.Response) CompetitionSeasonSummary {

	doc, err := goquery.NewDocumentFromReader(res.Body)
		
	if err != nil {
		log.Fatal("Failed to parse the HTML document", err)
	}

	rows := doc.Find("tr")

	var matchesData []MatchSummary

	for i, row := range rows.Nodes {

		if i == 0 {
			continue
		}

		// ignore row if it has the specific class "spacer partial_table result_all"
		// do this with brute force / ignore all rows with a class name
		if len(row.Attr) > 0 {
			continue
		}

		// iterate over each column, process according to its type, and add to a map representing the row
		matchMap := MatchSummary{data: make(map[string]string)}
		for column := range row.ChildNodes() {
			// ignore the column if it doesn't have enough attrs
			if len(column.Attr) < 1 {
				continue
			}
			parseColumn(column, matchMap)
		}
		matchesData = append(matchesData, matchMap)
	}
	return CompetitionSeasonSummary{Data: matchesData, CompetitionName: "Premier League", CompetitionSeason: "21-22"}
}

func parseColumn(column *html.Node, matches MatchSummary) {

	//map is passed by reference so we don't need to return anything.

	// work out what type of stat this particular column corresponds to
	var statType string
	if column.Attr[1].Key == "data-stat" {
		//the data-stat custom attr should nearly always be at index 1
		statType = column.Attr[1].Val
		// fmt.Printf("column.Attr[1].Val: %v\n", column.Attr[1].Val)
	} else {
		//the winning team has the data-stat as the next attr, as the previous will identify it for bold highlighting
		statType = column.Attr[2].Val
		// fmt.Printf("column.Attr[2].Val: %v\n", column.Attr[2].Val)
	}

	// process the column depending on the type of stat included
	switch statType {
	case "home_team", "away_team", "date":
		for child := range column.Descendants(){
			if len(child.Attr) > 0 {
				// fmt.Println(child.Attr[0].Val)
				matches.data[fmt.Sprintf("%s-url",statType)] = child.Attr[0].Val
			} else {
				matches.data[statType] = child.Data
				// fmt.Println(child.Data)
			}
		}
	case "score":
		for child := range column.Descendants(){
			if len(child.Attr) > 0 {
				matches.data["match-report"] = child.Attr[0].Val
				// fmt.Println(child.Attr[0].Val)
			} else {
				matches.data[statType] = child.Data
				// fmt.Println(child.Data)
			}
		}
	case "attendance", "venue", "referee", "dayofweek", "home_xg", "away_xg":
		for child := range column.Descendants(){
			matches.data[statType] = child.Data
		}
	case "start_time":
		if column.FirstChild != nil {
			if len(column.FirstChild.Attr) > 2 {
				matches.data["venue-time"] = column.FirstChild.Attr[3].Val
			}
			// fmt.Println(column.FirstChild)
		}
			// if len(column.Attr) > 2 {
		// 	matchMap["venue-time"] = column.FirstChild.Attr[3].Val
		// 	// fmt.Println(column.FirstChild.Attr[3].Val)
		// }
	default:
		// fmt.Println("no data")
		return
	}
}

func PrintMatches(matches []MatchSummary, limit int) {

	for i, match := range matches {
		PrintMatch(match, false)
		if i > limit {
			return
		}
	}
}

func PrintMatch(match MatchSummary, extended bool) {

	fmt.Println("=============================================")
	fmt.Printf("%s %s %s\n", match.data["home_team"], match.data["score"], match.data["away_team"])
	fmt.Printf("Time: %s, %s, %s\n", match.data["venue-time"], match.data["dayofweek"], match.data["date"])
	fmt.Printf("Referee: %s\n", match.data["referee"])
	fmt.Printf("Venue: %s, Attendance: %s\n", match.data["venue"], match.data["attendance"])
	fmt.Println("=============================================")

	if extended {
		for key, val := range match.data {
			fmt.Println(key, val)
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
}