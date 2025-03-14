package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/semidesnatada/fscraper/config"
	"github.com/semidesnatada/fscraper/database"

	_ "github.com/lib/pq"
)

func main() {
	
	const DB_URL = "postgres://seanlowery:@localhost:5432/fscraped?sslmode=disable"
	s := config.CreateState(DB_URL)


	// s.DeleteDatabases()	
	// client.ScrapeLeagues(&s)

	games, err := s.DB.GetGamesByTeamAndSeason(
		context.Background(),
		database.GetGamesByTeamAndSeasonParams{
			Name: "Liverpool",
			Season: "2015-2016",
			
		},
	)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("whoops - something wrong")
		os.Exit(1)
	}

	for i, game := range games {
		fmt.Println("=================================")
		fmt.Printf("Match %d\n", i+1)
		fmt.Printf("%s 	%d:%d 	%s\n",game.HomeTeam,game.HomeGoals,game.AwayGoals,game.AwayTeam)
		fmt.Printf("%v , 	%s\n",game.Date.Format(time.DateOnly),game.Stadium)
		fmt.Println("=================================")
	}


	fmt.Println("===================================================")
	fmt.Println()
	fmt.Println("Script concluded")
	fmt.Println()
	fmt.Println("===================================================")

	os.Exit(0)
}
