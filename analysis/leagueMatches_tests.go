package analysis

import (
	"context"
	"fmt"

	"github.com/semidesnatada/fscraper/client"
	"github.com/semidesnatada/fscraper/config"
	"github.com/semidesnatada/fscraper/database"
)

func CheckAllLeagueTables(s *config.State) error {
	// function which checks whether the leagues database is valid, by checking every league table query.

	comps, err := s.DB.GetUniqueCompetitionSeasons(context.Background())
	if err != nil {
		return err
	}

	// knockoutDetails := client.GetKnockoutParams()
	leagueDetails := client.GetLeagueParams()

	for _, comp := range comps {
		// only proceed if this 
		for _, league := range leagueDetails {
			if league.Name == comp.Name {
				checkErr := CheckLeagueTable(s, comp.Name, comp.Season)
				if checkErr != nil {
					return checkErr
				}
			}
		}
	}
	return nil
}

func CheckLeagueTable(s *config.State, leagueName string, leagueSeason string) error {
	// checks whether a given league table is valid mathematically. 
	// this does not check against a second source of truth for the matches played;
	// rather simply whether the stats aren't obviously wrong
	rows, err := s.DB.GetCompetitionTable(
		context.Background(),
		database.GetCompetitionTableParams{
			Name: leagueName,
			Season: leagueSeason,
		},
	)
	if err != nil {
		return err
	}

	noOfTeams := len(rows)
	var totalWins int
	var totalDraws int
	var totalLosses int
	var totalCalcPoints int
	var totalNetPoints int
	var totalGoalsScored int
	var totalGoalsConceded int
	var totalGamesPlayed int
	gpCheck := rows[0].GamesPlayed

	for _, row := range rows {
		if row.GamesPlayed != 2 * int64(noOfTeams - 1) {
			return fmt.Errorf("error with the games played by %s in the season %s %s",row.TeamName, leagueName, leagueSeason)
		}
		if row.GamesPlayed != gpCheck {
			return fmt.Errorf("error with the games played by %s in the season %s %s",row.TeamName, leagueName, leagueSeason)
		}
		if row.Points != (row.Draws + row.Wins * 3) {
			return fmt.Errorf("error with the points scored by %s in the season %s %s",row.TeamName, leagueName, leagueSeason)
		}
		if (row.GoalsScored - row.GoalsConceded) != row.GoalDifference {
			return fmt.Errorf("error with the goal difference by %s in the season %s %s",row.TeamName, leagueName, leagueSeason)
		}
		totalWins += int(row.Wins)
		totalDraws += int(row.Draws)
		totalLosses += int(row.Losses)
		totalCalcPoints += (3 * int(row.Wins) + int(row.Draws))
		totalNetPoints += int(row.Points)
		totalGoalsScored += int(row.GoalsScored)
		totalGoalsConceded += int(row.GoalsConceded)
		totalGamesPlayed += int(row.GamesPlayed)

	}

	if totalWins != totalLosses {
		return fmt.Errorf("issue with wins and losses in %s, %s",leagueName, leagueSeason)
	}
	if totalNetPoints != totalCalcPoints {
		return fmt.Errorf("issue with calculated points in %s, %s",leagueName, leagueSeason)
	}
	if totalGoalsScored != totalGoalsConceded {
		return fmt.Errorf("issue with goals scored / conceded in %s, %s",leagueName, leagueSeason)
	}
	if totalGamesPlayed != (noOfTeams - 1) * 2 * noOfTeams {
		return fmt.Errorf("issue with games played in %s, %s",leagueName, leagueSeason)
	}
	return nil
}