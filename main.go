package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/semidesnatada/fscraper/config"
	"github.com/semidesnatada/fscraper/database"

	_ "github.com/lib/pq"
)

func main() {
	
	const DB_URL = "postgres://seanlowery:@localhost:5432/fscraped?sslmode=disable"
	s := config.CreateState(DB_URL)

	seasonName := "Premier-League"
	seasonYear := "2023-2024"

	rows, err := s.DB.GetCompetitionTable(
		context.Background(),
		database.GetCompetitionTableParams{
			Name: seasonName,
			Season: seasonYear,
		},
	)
	if err != nil {
		fmt.Printf("error getting table %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println()
	fmt.Println("===========================================================================================")
	fmt.Printf("%s %s\n",seasonName, seasonYear)
	fmt.Println("===========================================================================================")
	fmt.Println("  Team  			P	W	D	L	GF	GA	GD	PTS")
	for place, row := range rows {
		var strPlace string
		if place < 9 {
			strPlace = strconv.Itoa(place + 1) + " " 
		} else {
			strPlace = strconv.Itoa(place + 1)
		}
		teamName := row.TeamName + strings.Repeat(" ", 20 - len(row.TeamName))

		fmt.Printf("%s %s       	%d	%d	%d	%d	%d	%d	%d	%d\n",
			strPlace, teamName, row.GamesPlayed, row.Wins, row.Draws, row.Losses, row.GoalsScored, row.GoalsConceded, row.GoalDifference, row.Points)
	}
	fmt.Println("===========================================================================================")
	fmt.Println()

	// s.DeleteDatabases()	
	// client.ScrapeLeagues(&s)

	// GetGamesTest(&s)


	fmt.Println("===================================================")
	fmt.Println()
	fmt.Println("Script concluded")
	fmt.Println()
	fmt.Println("===================================================")

	os.Exit(0)
}


func GetGamesTest(s *config.State) {
	games, err := s.DB.GetGamesByTeamAndSeason(
		context.Background(),
		database.GetGamesByTeamAndSeasonParams{
			Name: "Newcastle Utd",
			Season: "1996-1997",
			
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
}
