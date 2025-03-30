package analysis

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/semidesnatada/fscraper/config"
	"github.com/semidesnatada/fscraper/database"
)

func PrintAllLeagueTables(s *config.State) error {

	seasons, err := s.DB.GetUniqueCompetitionSeasons(context.Background())
	if err != nil {
		return err
	}

	for _, season := range seasons {

		indErr := GetAndPrintLeagueTable(s, season.Name, season.Season)
		if indErr != nil {
			return indErr
		}
	}
	return nil
}

func GetAndPrintClubCompetitionResultsTable(s *config.State, teamName string) error {
	rows, err := s.DB.GetAllClubCompetitionResults(context.Background(), teamName)
	if err != nil {
		return err
	}
	fmt.Println()
	fmt.Println("============================================================================================================")
	fmt.Printf("%s record in all competitions\n",teamName)
	fmt.Println("============================================================================================================")
	fmt.Println("  Competition  					P	W	D	L	GF	GA	GD	PTS")
	for comp, row := range rows {
		var strPlace string
		if comp < 9 {
			strPlace = strconv.Itoa(comp + 1) + " " 
		} else {
			strPlace = strconv.Itoa(comp + 1)
		}
		compName := fmt.Sprintf("%s %s%s", row.CompetitionName, row.CompetitionSeason, strings.Repeat(" ", 30 - len(row.CompetitionName)- len(row.CompetitionSeason)))

		fmt.Printf("%s %s       	%d	%d	%d	%d	%d	%d	%d	%d\n",
			strPlace, compName, row.GamesPlayed, row.Wins, row.Draws, row.Losses, row.GoalsScored, row.GoalsConceded, row.GoalDifference, row.Points)
	}
	fmt.Println("============================================================================================================")
	fmt.Println()
	return nil
}


func GetAndPrintAllTimeLeagueTable(s *config.State, seasonName string) error {
	rows, err := s.DB.GetAllTimeCompetitionTable(context.Background(),seasonName)
	if err != nil {
		return err
	}
	fmt.Println()
	fmt.Println("===========================================================================================")
	fmt.Printf("All Time %s\n",seasonName)
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
	return nil
}

func GetAndPrintLeagueTable(s *config.State, seasonName, seasonYear string) error {
	rows, err := s.DB.GetCompetitionTable(
		context.Background(),
		database.GetCompetitionTableParams{
			Name: seasonName,
			Season: seasonYear,
		},
	)
	if err != nil {
		return err
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
	return nil
}

func PrintScriptEnd() {
	fmt.Println("===========================================================================================")
	fmt.Println()
	fmt.Println("Script concluded")
	fmt.Println()
	fmt.Println("===========================================================================================")
}

func GetGamesTeamSeason(s *config.State, teamName, competition, season string) {
	games, err := s.DB.GetLeagueGamesByTeamAndSeason(
		context.Background(),
		database.GetLeagueGamesByTeamAndSeasonParams{
			Name: teamName,
			Name_2: competition,
			Season: season,
			
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
