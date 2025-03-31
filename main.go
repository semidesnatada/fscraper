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

	// testDBShortestPath(&s)
	// testPlayerPathDetailedStats(&s)
	// testGettingAllPaths(&s)

	// fmt.Println()
	// fmt.Println()
	// fmt.Println()

	analysis.TestPlayerMatchData(&s, 9)

	// this call now done via testing
	// analysis.CheckAllLeagueTables(&s)

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
	
	gFilename := "master_graph_dump_test"

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

func testPlayerPathDetailedStats(s *config.State) {

	// p1 := "/en/players/003cf4d1/Jayden-Danns"
	// p2 := "/en/players/0ac94a23/Denis-Cheryshev"
	// p2 := "/en/players/042cac2d/Moritz-Stoppelkamp"
	// p1 := "/en/players/03760df0/Emiliano-Viviano"
	// p2 := "/en/players/030928e4/Slaven-Bilic"
	// p2 := "/en/players/032f766b/Ali-Al-Habsi"
	// p2 := "/en/players/042e8a49/Kingsley-Coman"

	// p1 := "/en/players/dc001a06/Xabi-Alonso"
	// p2 := "/en/players/1ddbb0da/Iker-Casillas"

	// p1 := "/en/players/e06683ca/Virgil-van-Dijk"
	// p1 := "/en/players/cd1acf9d/Trent-Alexander-Arnold"
	// p2 := "/en/players/4c370d81/Roberto-Firmino"
	// p1 := "/en/players/c691bfe2/Sadio-Mane"
	// p2 := "/en/players/38c7feef/Alex-Oxlade-Chamberlain"

	p1 := "/en/players/002d06bb/Sylvain-Distin"
	p2 := "/en/players/08511d65/Sergio-Ramos"
	
	gFilename := "master_graph_dump_test"

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

	out, err := analysis.GetPathDetailedStatsFromUrls(s, path3)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	out.PrintPath(s)

}

func testGettingAllPaths(s *config.State) {

	//10097 : [/en/players/96cf7b61/Wes-Foderingham /en/players/57827369/Ben-Brereton /en/players/4bd414c1/Raul-Albiol /en/players/08511d65/Sergio-Ramos]
	// 86 : [/en/players/824c7343/Graham-Potter /en/players/10e485ae/Eyal-Berkovic /en/players/9be05f40/Steve-McManaman /en/players/1ddbb0da/Iker-Casillas /en/players/08511d65/Sergio-Ramos]
	// 87 : [/en/players/892d2b5f/Ali-Dia /en/players/10e485ae/Eyal-Berkovic /en/players/9be05f40/Steve-McManaman /en/players/1ddbb0da/Iker-Casillas /en/players/08511d65/Sergio-Ramos]
	/// en/players/231f7c03/Glenn-Hoddle /en/players/8cd0120b/Scott-Minto /en/players/f1e8372d/Jermain-Defoe /en/players/21a66f6a/Harry-Kane /en/players/042e8a49/Kingsley-Coman
	
	gFilename := "master_graph_dump_test"

	Q, e2 := analysis.ReadGraphFromFile(gFilename)
	if e2 != nil {
		fmt.Println(e2.Error())
		os.Exit(1)
	}

	// p2 := "/en/players/08511d65/Sergio-Ramos"
	p2 := "/en/players/042e8a49/Kingsley-Coman"
	paths, err := Q.GetPathsBelowGivenLength(p2, 5, s)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for i, path := range paths {
		if i == 4{
			fmt.Printf("All players a distance %d away from %s\n", i, p2)
			for j, x := range path {
				fmt.Println(j, ":", x)
			}}
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