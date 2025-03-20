package main

import (
	"fmt"
	"os"

	"github.com/semidesnatada/fscraper/analysis"
	"github.com/semidesnatada/fscraper/client"
	"github.com/semidesnatada/fscraper/config"

	_ "github.com/lib/pq"
)

func main() {
	
// 8f1efd14-48af-451e-b99c-eb52fe29f88b == charlie hartfield for QPR


	const DB_URL = "postgres://seanlowery:@localhost:5432/fscraped?sslmode=disable"
	s := config.CreateState(DB_URL)

	s.DeleteAllDatabases()
	// s.DeleteDetailedDatabases()
	// s.DeleteSummaryDatabases()	

	client.ScrapeLeagues(&s)
	client.ScrapeKnockouts(&s)

	// client.DetailedMatchScraper(&s)

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