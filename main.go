package main

import (
	"fmt"
	"os"

	"github.com/semidesnatada/fscraper/analysis"
	"github.com/semidesnatada/fscraper/config"

	_ "github.com/lib/pq"
)

func main() {
	
// 8f1efd14-48af-451e-b99c-eb52fe29f88b == charlie hartfield for QPR

// Check these:
// = 1426 Started scraping url: /en/matches/4441ad45/Real-Madrid-Compostela-March-22-1998-La-Liga 
// == Already scraped, moving on
// = 1427 Started scraping url: /en/matches/4441c868/Leeds-United-Southampton-April-3-1996-Premier-League 
// == Already scraped, moving on
// = 1428 Started scraping url: /en/matches/4442b38e/Eintracht-Frankfurt-Augsburg-December-7-2024-Bundesliga 
// == Successfully finished scraping url: /en/matches/4442b38e/Eintracht-Frankfurt-Augsburg-December-7-2024-Bundesliga
// === Started storing url: /en/matches/4442b38e/Eintracht-Frankfurt-Augsburg-December-7-2024-Bundesliga 
// ==== Completed storing url: /en/matches/4442b38e/Eintracht-Frankfurt-Augsburg-December-7-2024-Bundesliga 



	const DB_URL = "postgres://seanlowery:@localhost:5432/fscraped?sslmode=disable"
	s := config.CreateState(DB_URL)

	// s.DeleteAllDatabases()
	// s.DeleteDetailedDatabases()
	// s.DeleteSummaryDatabases()	

	// client.ScrapeLeagues(&s)
	// client.ScrapeKnockouts(&s)

	// client.DetailedMatchScraper(&s)
	
	// xErr := client.DeleteAndRescrapeLeagueMatches(&s)
	// if xErr != nil {
	// 	fmt.Println(xErr.Error())
	// }

	testDBShortestPath(&s)

	// analysis.TestPlayerMatchData(&s, 9)
	// analysis.PrintAllLeagueTables(&s)

	// checkErr := analysis.CheckAllLeagueTables(&s)
	// if checkErr != nil {
	// 	fmt.Println()
	// 	fmt.Printf("error checking league tables %w\n", checkErr)
	// 	fmt.Println()
	// 	os.Exit(1)
	// }

	// err := analysis.GetAndPrintKnockoutDraw(&s, "Champions-League", "2023-2024", 3)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(1)
	// }

	// tester := client.CompetitionSeasonSummary{
	// 	CompetitionName: "Championship",
	// 	CompetitionSeason: "2016-2017",
	// 	CompetitionOnlineID: "10",
	// 	Url: "https://fbref.com/en/comps/10/2016-2017/schedule/2016-2017-Championship-Scores-and-Fixtures",
	// }

	// tester2 := client.CompetitionSeasonSummary{
	// 	CompetitionName: "Champions-League",
	// 	CompetitionSeason: "2023-2024",
	// 	CompetitionOnlineID: "8",
	// 	Url: "https://fbref.com/en/comps/8/2023-2024/schedule/2023-2024-Champions-League-Scores-and-Fixtures",
	// }

	// client.ScrapeLeagueFromUrl(&tester2)

	// teamName := "Newcastle Utd"
	// season := "1996-1997"
	// analysis.GetGamesTeamSeason(&s, teamName, season)

	// TestLeagueTableQuery(&s)
	// TestGamesStorage(&s)

	// analysis.GetGamesTeamSeason(&s, "Hamburger SV", "Bundesliga", "2014-2015")
	// err := analysis.PrintAllLeagueTables(&s)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(1)
	// }

	// err := analysis.GetAndPrintAllTimeLeagueTable(&s, "Ligue-1")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(1)
	// }

	// err := analysis.GetAndPrintClubCompetitionResultsTable(&s, "Norwich City")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(1)
	// }

	// err := analysis.PrintAllLeagueTables(&s)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(1)
	// }

	analysis.PrintScriptEnd()
	os.Exit(0)
}

func TestLeagueTableQuery(s *config.State) {
	seasonName := "Championship"
	seasonYear := "2015-2016"

	err := analysis.GetAndPrintLeagueTable(s, seasonName, seasonYear)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func testDBShortestPath(s *config.State) {
	g, err := analysis.NewGraphFromDB(s)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	p1 := "/en/players/002d06bb/Sylvain-Distin"
	p2 := "/en/players/08511d65/Sergio-Ramos"

	path, sharedmins, err2 := g.GetShortestConnectionBetweenPlayerUrls(s, p1, p2)
	if err2 != nil {
		fmt.Println(err2.Error())
		os.Exit(1)
	}
	fmt.Printf("shortest path between %s and %s is:\n", p1, p2)
	for i, step := range path {
		if i < len(path) -1 {
		fmt.Printf("%d : %s, played %d minutes with:\n", i, step, sharedmins[i+1])
	} else {
		fmt.Printf("%d : %s\n", i, step)
	}
	}
	
	gFilename := "test_graph_dump"

	e:= g.Write(gFilename)
	if e != nil {
		fmt.Println(e.Error())
	}

	Q, e2 := analysis.ReadGraphFromFile(gFilename)
	if e2 != nil {
		fmt.Println(e2.Error())
		os.Exit(1)
	}

	path3, sharedmins3, err4 := Q.GetShortestConnectionBetweenPlayerUrls(s, p1, p2)
	if err4 != nil {
		fmt.Println(err4.Error())
		os.Exit(1)
	}
	fmt.Printf("shortest path between %s and %s is:\n", p1, p2)
	for i, step := range path3 {
		if i < len(path3) -1 {
		fmt.Printf("%d : %s, played %d minutes with:\n", i, step, sharedmins3[i+1])
	} else {
		fmt.Printf("%d : %s\n", i, step)
	}
	}
}

// func TestGamesStorage(s *config.State) {
// 	data, err := s.DB.GetGamesByTeamAndSeason(
// 		context.Background(),
// 		database.GetGamesByTeamAndSeasonParams{
// 			Name:"Hamburger SV",
// 			Name_2:"Bundesliga",
// 			Season:"2014-2015",
// 		},
// 	)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		os.Exit(1)
// 	}
// 	for _, row := range data {
// 		fmt.Println(row.Date, row.HomeTeam, row.HomeGoals, row.AwayGoals, row.AwayTeam, row.Stadium)
// 	}
// }