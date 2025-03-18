package analysis

import (
	"context"
	"fmt"

	"github.com/semidesnatada/fscraper/config"
	"github.com/semidesnatada/fscraper/database"
)


func GetAndPrintKnockoutDraw(s *config.State, compName, compSeason string, roundsLimit int) error {

	// do this once tested db functionality. easiest with SQL queries.
	//e.g. query to find out number of rounds and number of games in each round
	//then , query he matches in the round with smallest number of matches (final) and so on, up to roundsLimit
	//then print the matches ordered based on rounds

	roundMatchPairs, xErr := s.DB.GetMatchesInEachRoundForGivenComp(
		context.Background(),
		database.GetMatchesInEachRoundForGivenCompParams{
			Name: compName,
			Season: compSeason,
		},
	)
	if xErr != nil {
		return xErr
	}
	fmt.Println()
	fmt.Println(roundMatchPairs)
	fmt.Println()
	// for _, y := range x {
	// 	fmt.Println(y)
	// }

	matchesInRound, err := s.DB.GetKnockoutGamesByRoundAndSeason(
		context.Background(),
		database.GetKnockoutGamesByRoundAndSeasonParams{
			Round: "Semi-finals",
			Name: compName,
			Season: compSeason,
		},
	)
	if err != nil {
		return nil
	}
	// to be used if there are multiple legs to consider
	handleTwoLeggedRound(matchesInRound)


	return nil

}

func handleTwoLeggedRound(matches []database.GetKnockoutGamesByRoundAndSeasonRow) {

	secLegMap := map[int]int{}

	for i, game_1 := range matches {
		for j, game_2 := range matches {
			if game_1.HomeTeam == game_2.AwayTeam {
				if _, ok := secLegMap[j]; !ok {
					secLegMap[i] = j
				}
			}
		}
	}

	for firstLeg, secondLeg := range secLegMap {
		printTwoLeggedKnockoutMatchSummary(matches[firstLeg], matches[secondLeg])
		fmt.Println()
	}
}


func printTwoLeggedKnockoutMatchSummary(game_1, game_2 database.GetKnockoutGamesByRoundAndSeasonRow) {
	
	if !game_1.WentToPens && !game_2.WentToPens {
		fmt.Printf(">> %s %d:%d %s\n", game_1.HomeTeam, game_1.HomeGoals, game_1.AwayGoals, game_1.AwayTeam)
		fmt.Printf(">> %s %d:%d %s\n", game_2.HomeTeam, game_2.HomeGoals, game_2.AwayGoals, game_2.AwayTeam)
		homeTot := game_1.HomeGoals + game_2.AwayGoals
		awayTot :=  game_2.HomeGoals + game_1.AwayGoals
		if homeTot > awayTot {
			fmt.Printf(">>>> %s win %d:%d on aggregate\n", game_1.HomeTeam, homeTot, awayTot)
		} else if game_1.HomeGoals + game_2.AwayGoals < game_2.HomeGoals + game_1.AwayGoals {
			fmt.Printf(">>>> %s win %d:%d on aggregate\n", game_2.HomeTeam, homeTot, awayTot)
		} else {
			fmt.Printf(">>>> aggregate score %d:%d (work out who won based on who goes through to next round).\n", homeTot, awayTot)
		}
	} else if game_1.WentToPens{
		fmt.Printf(">> %s %d:%d %s\n", game_2.HomeTeam, game_2.HomeGoals, game_2.AwayGoals, game_2.AwayTeam)
		fmt.Printf(">> %s %d:%d %s\n", game_1.HomeTeam, game_1.HomeGoals, game_1.AwayGoals, game_1.AwayTeam)
		fmt.Printf(">>> pens: %s %d:%d %s\n", game_1.HomeTeam, game_1.HomePens.Int32, game_1.AwayPens.Int32, game_1.AwayTeam)
		if game_1.HomePens.Int32 > game_1.AwayPens.Int32 {
			fmt.Printf(">>>> %s win after penalties\n", game_1.HomeTeam)
		} else {
			fmt.Printf(">>>> %s win after penalties\n", game_1.AwayTeam)
		}
	} else {
		fmt.Printf(">> %s %d:%d %s\n", game_1.HomeTeam, game_1.HomeGoals, game_1.AwayGoals, game_1.AwayTeam)
		fmt.Printf(">> %s %d:%d %s\n", game_2.HomeTeam, game_2.HomeGoals, game_2.AwayGoals, game_2.AwayTeam)
		fmt.Printf(">>> pens: %s %d:%d %s\n", game_2.HomeTeam, game_2.HomePens.Int32, game_2.AwayPens.Int32, game_2.AwayTeam)
		if game_2.HomePens.Int32 > game_2.AwayPens.Int32 {
			fmt.Printf(">>>> %s win after penalties\n", game_2.HomeTeam)
		} else {
			fmt.Printf(">>>> %s win after penalties\n", game_2.AwayTeam)
		}
	}
	
}