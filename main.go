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
	
	const DB_URL = "postgres://seanlowery:@localhost:5432/fscraped?sslmode=disable"
	// s := config.CreateState(DB_URL)

	// s.DeleteDatabases()	
	// client.ScrapeLeagues(&s)

	tester := client.CompetitionSeasonSummary{
		CompetitionName: "Championship",
		CompetitionSeason: "2016-2017",
		CompetitionOnlineID: "10",
		Url: "https://fbref.com/en/comps/10/2016-2017/schedule/2016-2017-Championship-Scores-and-Fixtures",
	}

	client.ScrapeLeagueFromUrl(&tester)

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

	// err := analysis.GetAndPrintAllTimeLeagueTable(&s, "La-Liga")
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