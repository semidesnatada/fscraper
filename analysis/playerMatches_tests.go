package analysis

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/semidesnatada/fscraper/client"
	"github.com/semidesnatada/fscraper/config"
	"github.com/semidesnatada/fscraper/database"
)

// func which checks whether the table prepared using the leaguematches database yields an equal result to the player matches database.

type MatchSummary struct {
	HomeTeam, AwayTeam string
	HomeGoals, AwayGoals int
}

func VerifyBothDBs(s *config.State) error {

	leagues := client.GenerateLeaguesForSearching()

	for _, league := range leagues {
		name := league.CompetitionName
		season := league.CompetitionSeason

		fmt.Printf("now verifying %s %s\n", name, season)

		matchesCheck := VerifyPlayerMatchesLeagueMatchesDBs(s, name, season)
		tableCheck := VerifyPlayerMatchesLeagueTableDBs(s, name, season)

		if matchesCheck != nil || tableCheck != nil {
			return fmt.Errorf("data in the following comp is incorrect %s %s.\n Matches error message: %w ;\n Table error message: %w",
								name, season, matchesCheck, tableCheck)
			// fmt.Printf("data in the following comp is incorrect %s %s.\n Matches error message: %v ;\n Table error message: %v\n",
			// 					name, season, matchesCheck, tableCheck)
		}
	}
	return nil
}

func VerifyPlayerMatchesLeagueMatchesDBs(s *config.State, leagueName, leagueSeason string) error {
	
	gamesLM, err := s.DB.GetLeagueGamesBySeason(
		context.Background(),
		database.GetLeagueGamesBySeasonParams{
			Name: leagueName,
			Season: leagueSeason,
		},
	)
	if err != nil {
		return err
	}

	lmMap := make(map[uuid.UUID]MatchSummary)

	for _, match := range gamesLM {
		lmMap[match.ID] = MatchSummary{
			HomeTeam: match.HomeTeam,
			AwayTeam: match.AwayTeam,
			HomeGoals: int(match.HomeGoals),
			AwayGoals: int(match.AwayGoals),
		}
	}

	gamesPM, err2 := s.DB.GetMatchRecordsFromPM(
		context.Background(),
		database.GetMatchRecordsFromPMParams{
			Name: leagueName,
			Season: leagueSeason,
		},
	)
	if err2 != nil {
		return err2
	}

	for _, game := range gamesPM {
		comparator := lmMap[game.MatchID]
		homeGls := game.HomeGoalsScored + game.AwayOgs
		awayGls := game.AwayGoalsScored + game.HomeOgs
		if game.HomeTeamName != comparator.HomeTeam || game.AwayTeamName != comparator.AwayTeam {
			return fmt.Errorf("issue with team names not correctly identified, %s : %s | home team names: %s ; %s | away team names: %s ; %s", leagueName, leagueSeason, game.HomeTeamName, comparator.HomeTeam, game.AwayTeamName, comparator.AwayTeam)
			// fmt.Printf("issue with team names not correctly identified, %s : %s | home team names: %s ; %s | away team names: %s ; %s\n", leagueName, leagueSeason, game.HomeTeamName, comparator.HomeTeam, game.AwayTeamName, comparator.AwayTeam)
		}
		if int(homeGls) != comparator.HomeGoals || int(awayGls) != comparator.AwayGoals {
			return fmt.Errorf("issue with goals scored in match: %s v %s | Score is recorded in PM as %d : %d, when it should be %d : %d", game.HomeTeamName, game.AwayTeamName, game.HomeGoalsScored, game.AwayGoalsScored, comparator.HomeGoals, comparator.AwayGoals)
			// fmt.Printf("issue with goals scored in match: %s v %s | Score is recorded in PM as %d : %d, when it should be %d : %d\n", game.HomeTeamName, game.AwayTeamName, game.HomeGoalsScored, game.AwayGoalsScored, comparator.HomeGoals, comparator.AwayGoals)
		}
	}
	return nil
}

func PrintLeagueTableFromBothDBs(s *config.State, leagueName, leagueSeason string) error {
	// this is the league matches table
	rowsLM, errLM := s.DB.GetCompetitionTable(
		context.Background(),
		database.GetCompetitionTableParams{
			Name: leagueName,
			Season: leagueSeason,
		},
	)
	if errLM != nil {
		return errLM
	}
	// this is the player matches table
	rowsPM, errPM := s.DB.GetCompTableFromPM(
		context.Background(),
		database.GetCompTableFromPMParams{
			Name: leagueName,
			Season: leagueSeason,
		},
	)
	if errPM != nil {
		return errPM
	}
	fmt.Println("League Matches Table")
	PrintLMTable(rowsLM, leagueName, leagueSeason)
	fmt.Println("Player Matches Table")
	PrintPMTable(rowsPM, leagueName, leagueSeason)
	return nil
}

func VerifyPlayerMatchesLeagueTableDBs(s *config.State, leagueName, leagueSeason string) error {
	// this is the league matches table
	rowsLM, errLM := s.DB.GetCompetitionTable(
		context.Background(),
		database.GetCompetitionTableParams{
			Name: leagueName,
			Season: leagueSeason,
		},
	)
	if errLM != nil {
		return errLM
	}
	// this is the player matches table
	rowsPM, errPM := s.DB.GetCompTableFromPM(
		context.Background(),
		database.GetCompTableFromPMParams{
			Name: leagueName,
			Season: leagueSeason,
		},
	)
	if errPM != nil {
		return errPM
	}

	for i, row := range rowsLM {
		if row.TeamName != rowsPM[i].TeamName {
			// fmt.Println("error with ")
			return fmt.Errorf("misaligned team names for %s : %s | %s vs %s", leagueName, leagueSeason, row.TeamName, rowsPM[i].TeamName)
		}
		if row.GamesPlayed != rowsPM[i].MatchesPlayed {
			// fmt.Println("error with ")
			return fmt.Errorf("misaligned matches played for %s : %s , %s | %d vs %d", leagueName, leagueSeason, row.TeamName, row.GamesPlayed, rowsPM[i].MatchesPlayed)
		}
		if row.Wins != rowsPM[i].Wins {
			// fmt.Println("error with ")
			return fmt.Errorf("misaligned win numbers for %s : %s , %s | %d vs %d", leagueName, leagueSeason, row.TeamName, row.Wins, rowsPM[i].Wins)
		}
		if row.Draws != rowsPM[i].Draws {
			// fmt.Println("error with ")
			return fmt.Errorf("misaligned draw numbers for %s : %s , %s | %d vs %d", leagueName, leagueSeason, row.TeamName, row.Draws, rowsPM[i].Draws)
		}
		if row.Losses != rowsPM[i].Losses {
			// fmt.Println("error with ")
			return fmt.Errorf("misaligned losses numbers for %s : %s , %s | %d vs %d", leagueName, leagueSeason, row.TeamName, row.Losses, rowsPM[i].Losses)
		}
		if row.GoalsScored != rowsPM[i].GoalsScored {
			// fmt.Println("error with ")
			return fmt.Errorf("misaligned goals scored numbers for %s : %s , %s | %d vs %d", leagueName, leagueSeason, row.TeamName, row.GoalsScored, rowsPM[i].GoalsScored)
		}
		if row.GoalsConceded != rowsPM[i].GoalsConceded {
			// fmt.Println("error with ")
			return fmt.Errorf("misaligned goals conceded numbers for %s : %s , %s | %d vs %d", leagueName, leagueSeason, row.TeamName, row.GoalsConceded, rowsPM[i].GoalsConceded)
		}
		if row.GoalDifference != rowsPM[i].GoalDifference {
			// fmt.Println("error with ")
			return fmt.Errorf("misaligned goal difference numbers for %s : %s , %s | %d vs %d", leagueName, leagueSeason, row.TeamName, row.GoalsConceded, rowsPM[i].GoalsConceded)
		}
		if row.Points != rowsPM[i].Points {
			// fmt.Println("error with ")
			return fmt.Errorf("misaligned points numbers for %s : %s , %s | %d vs %d", leagueName, leagueSeason, row.TeamName, row.Points, rowsPM[i].Points)
		}
	}
	return nil
}

func PrintPMTable(t []database.GetCompTableFromPMRow, league, season string) {
	fmt.Println()
	fmt.Println("===========================================================================================")
	fmt.Printf("%s %s League table (Player Matches source)\n", league, season)
	fmt.Println("===========================================================================================")
	fmt.Printf("   %s	%s	%s	%s	%s	%s	%s	%s	%s\n",
				pSTS("Team", 25),
				pSTS("P", 4),
				pSTS("W", 4),
				pSTS("D",4),
				pSTS("L",4),
				pSTS("GF", 4),
				pSTS("GA", 4),
				pSTS("GD",4),
				pSTS("PTS",4),
	)

	for i, row := range t {
		fmt.Printf("%s%s	%s	%s	%s	%s	%s	%s	%s	%s\n",
		pNTS(i+1, 3),
		pSTS(row.TeamName, 25),
		pNTS(int(row.MatchesPlayed), 4),
		pNTS(int(row.Wins), 4),
		pNTS(int(row.Draws),4),
		pNTS(int(row.Losses),4),
		pNTS(int(row.GoalsScored), 4),
		pNTS(int(row.GoalsConceded), 4),
		pNTS(int(row.GoalDifference),4),
		pNTS(int(row.Points),4),
		)
	}
	fmt.Println("===========================================================================================")
	fmt.Println()

}

func PrintLMTable(t []database.GetCompetitionTableRow, league, season string) {
	fmt.Println()
	fmt.Println("===========================================================================================")
	fmt.Printf("%s %s League table (League Matches source)\n",league, season)
	fmt.Println("===========================================================================================")
	fmt.Printf("   %s	%s	%s	%s	%s	%s	%s	%s	%s\n",
				pSTS("Team", 25),
				pSTS("P", 4),
				pSTS("W", 4),
				pSTS("D",4),
				pSTS("L",4),
				pSTS("GF", 4),
				pSTS("GA", 4),
				pSTS("GD",4),
				pSTS("PTS",4),
	)

	for i, row := range t {
		fmt.Printf("%s%s	%s	%s	%s	%s	%s	%s	%s	%s\n",
		pNTS(i+1, 3),
		pSTS(row.TeamName, 25),
		pNTS(int(row.GamesPlayed), 4),
		pNTS(int(row.Wins), 4),
		pNTS(int(row.Draws),4),
		pNTS(int(row.Losses),4),
		pNTS(int(row.GoalsScored), 4),
		pNTS(int(row.GoalsConceded), 4),
		pNTS(int(row.GoalDifference),4),
		pNTS(int(row.Points),4),
		)
	}
	fmt.Println("===========================================================================================")
	fmt.Println()
}