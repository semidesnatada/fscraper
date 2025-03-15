package main

import (
	"fmt"
	"os"

	"github.com/semidesnatada/fscraper/analysis"
	"github.com/semidesnatada/fscraper/config"

	_ "github.com/lib/pq"
)

func main() {
	
	const DB_URL = "postgres://seanlowery:@localhost:5432/fscraped?sslmode=disable"
	s := config.CreateState(DB_URL)

	// s.DeleteDatabases()	
	// client.ScrapeLeagues(&s)

	// teamName := "Newcastle Utd"
	// season := "1996-1997"
	// analysis.GetGamesTeamSeason(&s, teamName, season)

	// TestLeagueTableQuery(&s)
	// err := analysis.PrintAllLeagueTables(&s)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(1)
	// }

	err := analysis.GetAndPrintAllTimeLeagueTable(&s, "La-Liga")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}


	// err := analysis.PrintAllLeagueTables(&s)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(1)
	// }

	analysis.PrintScriptEnd()
	os.Exit(0)
}

func TestLeagueTableQuery(s *config.State) {
	seasonName := "Premier-League"
	seasonYear := "2015-2016"

	err := analysis.GetAndPrintLeagueTable(s, seasonName, seasonYear)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
