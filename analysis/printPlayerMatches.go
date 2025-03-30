package analysis

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/semidesnatada/fscraper/config"
)

func TestPlayerMatchData(s *config.State, limit int) {

	// test_1(s, limit)
	// test_2(s, limit)
	test_3(s, limit*3)
	// test_4(s, limit)
	// test_5(s, limit)
	// test_6(s)

}

func padStat(i int32) string {

	if i < 10 {
		return fmt.Sprintf(" %d",i)
	} else {
		return fmt.Sprintf("%d", i)
	}

}

func test_1(s *config.State, limit int) {
	fmt.Println("===========================================================================================")
	fmt.Println()
	fmt.Println("Beginning Test 1")
	fmt.Println()
	test_1, err_1 := s.DB.GetNumberOfDistinctPlayersFieldedInLeagueByTeam(context.Background())
	if err_1 != nil {
		fmt.Println(err_1.Error())
		os.Exit(1)
	}
	fmt.Println("      : Competition		; 	Team				;Matches Played	;Players Fielded	")
	for i, item := range test_1 {
		var ind string
		if i < 9 {
			ind = fmt.Sprintf(" %d",i+1)
		} else {
			ind = fmt.Sprintf("%d",i+1)
		}
		compName := fmt.Sprintf("%s %s", item.CompetitionName, strings.Repeat(" ", 20 - len(item.CompetitionName)))
		teamName := fmt.Sprintf("%s %s", item.TeamName, strings.Repeat(" ", 20 - len(item.TeamName)))
		fmt.Printf("    %s: %s	;	%s	;	%d	;	%d\n",
		ind, compName, teamName, item.MatchesPlayed, item.PlayersFielded )
		if i == limit {
			break
		}
	}
	fmt.Println()
	fmt.Println("Test 1 Complete")
	fmt.Println()
	fmt.Println("===========================================================================================")
}

func test_2(s *config.State, limit int) {
	fmt.Println("===========================================================================================")
	fmt.Println()
	fmt.Println("Beginning Test 2")
	test_2, err_2 := s.DB.GetNumberOfGoalsScoredInEachLeagueSeasonByTeam(context.Background())
	if err_2 != nil {
		fmt.Println(err_2.Error())
		os.Exit(1)
	}
	fmt.Println("      : Competition			; 	Team			;Matches Played	;Goals Scored	;Players Fielded	")
	for i, item := range test_2 {
		var ind string
		if i < 9 {
			ind = fmt.Sprintf(" %d",i+1)
		} else {
			ind = fmt.Sprintf("%d",i+1)
		}
		compName := fmt.Sprintf("%s %s%s", item.CompetitionName, item.CompetitionSeason, strings.Repeat(" ", 30 - len(item.CompetitionName)- len(item.CompetitionSeason)))
		teamName := fmt.Sprintf("%s %s", item.TeamName,strings.Repeat(" ", 20 - len(item.TeamName)))

		fmt.Printf("    %s: %s	; 	%s	;	%d	;	%d	;	%d	\n",
		ind, compName, teamName, item.MatchesPlayed,item.TotalGoalsScored, item.PlayersFielded)
		if i == limit {
			break
		}
	}
	fmt.Println()
	fmt.Println("Test 2 Complete")
	fmt.Println()
	fmt.Println("===========================================================================================")
}

func test_3(s *config.State, limit int) {
	fmt.Println("===========================================================================================")
	fmt.Println()
	fmt.Println("Beginning Test 3a")
	test_3a, err_3a := s.DB.GetLeagueAllTimeTopScorers(context.Background())
	if err_3a != nil {
		fmt.Println(err_3a.Error())
		os.Exit(1)
	}
	fmt.Println("      : Player			; 	Competition			;Matches Played	;Goals Scored	")
	for i, item := range test_3a {
		var ind string
		if i < 9 {
			ind = fmt.Sprintf(" %d",i+1)
		} else {
			ind = fmt.Sprintf("%d",i+1)
		}
		compName := fmt.Sprintf("%s %s", item.CompetitionName,  strings.Repeat(" ", 30 - len(item.CompetitionName)))
		playerName := fmt.Sprintf("%s %s", item.PlayerName,strings.Repeat(" ", 30 - len(item.PlayerName)))

		fmt.Printf("    %s: %s	; 	%s	;	%d	;	%d	\n",
		ind,  playerName,compName, item.MatchesPlayed, item.TotalGoals)
		if i == limit {
			break
		}
	}
	fmt.Println()
	fmt.Println("Test 3a Complete")
	fmt.Println()
	fmt.Println("===========================================================================================")

	fmt.Println("===========================================================================================")
	fmt.Println()
	fmt.Println("Beginning Test 3b")
	test_3b, err_3b := s.DB.GetPlayerLeagueStats(context.Background())
	if err_3b != nil {
		fmt.Println(err_3b.Error())
		os.Exit(1)
	}
	fmt.Println("      : Player			; 	Competition			;Matches Played	;Goals Scored	")
	for i, item := range test_3b {
		var ind string
		if i < 9 {
			ind = fmt.Sprintf(" %d",i+1)
		} else {
			ind = fmt.Sprintf("%d",i+1)
		}
		compName := fmt.Sprintf("%s %s%s", item.CompetitionName, item.CompetitionSeason, strings.Repeat(" ", 30 - len(item.CompetitionName)- len(item.CompetitionSeason)))
		playerName := fmt.Sprintf("%s %s", item.PlayerName,strings.Repeat(" ", 30 - len(item.PlayerName)))

		fmt.Printf("    %s: %s	; 	%s	;	%d	;	%d	\n",
		ind,  playerName,compName, item.MatchesPlayed,item.TotalGoals)
		if i == limit {
			break
		}
	}
	fmt.Println()
	fmt.Println("Test 3b Complete")
	fmt.Println()
	fmt.Println("===========================================================================================")

	fmt.Println("===========================================================================================")
	fmt.Println()
	fmt.Println("Beginning Test 3c")
	test_3c, err_3c := s.DB.GetAllTimeTopScorers(context.Background())
	if err_3c != nil {
		fmt.Println(err_3c.Error())
		os.Exit(1)
	}
	fmt.Println("      : Player				;Matches Played	;Goals Scored	")
	for i, item := range test_3c {
		var ind string
		if i < 9 {
			ind = fmt.Sprintf(" %d",i+1)
		} else {
			ind = fmt.Sprintf("%d",i+1)
		}
		playerName := fmt.Sprintf("%s %s", item.PlayerName,strings.Repeat(" ", 30 - len(item.PlayerName)))

		fmt.Printf("    %s: %s	;	%d	;	%d	\n",
		ind,  playerName, item.MatchesPlayed,item.TotalGoals)
		if i == limit {
			break
		}
	}
	fmt.Println()
	fmt.Println("Test 3c Complete")
	fmt.Println()
	fmt.Println("===========================================================================================")
}

func test_4(s *config.State, limit int) {
	fmt.Println("===========================================================================================")
	fmt.Println("Beginning Test 4")

	// testUrl := "/en/players/c596fcb0/Callum-Wilson"
	testUrl := "/en/players/d70ce98e/Lionel-Messi"
	// testUrl := "/en/players/2c6835e5/Lewis-Miley"
	// testUrl := "/en/players/2b81295d/Raul"

	test_pre4, err_pre4 := s.DB.GetStatsForPlayerUrl(context.Background(), testUrl)
	if err_pre4 != nil {
		fmt.Println(err_pre4.Error())
		os.Exit(1)
	}
	test_4, err_4 := s.DB.GetPlayersPlayedWithByUrl(context.Background(), testUrl)
	if err_4 != nil {
		fmt.Println(err_4.Error())
		os.Exit(1)
	}
	fmt.Println()
	fmt.Println("Player stats are:")
	fmt.Printf("Name: %s , calculated mins played: %d\n", test_pre4.PlayerName, test_4[0].TotalMinsPlayed)

	fmt.Println()
	fmt.Println("Stats calculated with independent method:")
	fmt.Println("Matches played :", test_pre4.MatchesPlayed)
	fmt.Println("Total mins played :", test_pre4.TotalMinsPlayed)
	fmt.Println("Goals scorewd :", test_pre4.TotalGoals)
	fmt.Println("Pens scored :", test_pre4.TotalPens)
	fmt.Println("Own goals scored :", test_pre4.TotalOgs)
	fmt.Println("Total yellow cards :", test_pre4.TotalYellowCard)
	fmt.Println("Total red cards :", test_pre4.TotalRedCard)
	
	fmt.Println()
	fmt.Println("      : Colleague 				; Mins played together")
	for i, item := range test_4 {
		// if i == 0 {
		// 	continue
		// }
		var ind string
		if i < 9 {
			ind = fmt.Sprintf(" %d",i)
		} else {
			ind = fmt.Sprintf("%d",i)
		}
		cName := fmt.Sprintf("%s %s", item.ColleagueName, strings.Repeat(" ", 30 - len(item.ColleagueName)))
		fmt.Printf("    %s: %s		: %d\n",
		ind,  cName, item.TotalMinsPlayed)
		if i == limit {
			break
		}
	}
	fmt.Println()
	fmt.Println("Test 4 Complete")
	fmt.Println()
	fmt.Println("===========================================================================================")
}

func test_5(s *config.State, limit int) {
	fmt.Println("===========================================================================================")
	fmt.Println()
	fmt.Println("Beginning Test 5")
	test_5, err_5 := s.DB.GetMatchesWhereMinsDontAddUp(context.Background())
	if err_5 != nil {
		fmt.Println(err_5.Error())
		os.Exit(1)
	}
	fmt.Println("      : Total Mins			; 	Match Url")
	for i, item := range test_5 {
		var ind string
		if i < 9 {
			ind = fmt.Sprintf(" %d",i+1)
		} else {
			ind = fmt.Sprintf("%d",i+1)
		}

		fmt.Printf("    %s: %d	; 	%s	\n",
		ind, item.SquadMins, item.MatchUrl)
		if i == limit {
			break
		}
	}
	fmt.Println()
	fmt.Println("Test 5 Complete")
	fmt.Println()
	fmt.Println("===========================================================================================")
}

func test_6(s *config.State) {
	fmt.Println("===========================================================================================")
	fmt.Println()
	fmt.Println("Beginning Test 6")

	// testMatchUrl := "/en/matches/75e1f5c4/Valencia-Mallorca-November-28-2004-La-Liga"
	// testMatchUrl := "/en/matches/198facdd/Stuttgart-Wolfsburg-September-16-2017-Bundesliga"
	// testMatchUrl := "/en/matches/64b77eb3/Barcelona-Deportivo-La-Coruna-March-18-2000-La-Liga"
	testMatchUrl := "/en/matches/474ce5b7/Schalke-04-Bayer-Leverkusen-April-9-1996-Bundesliga"
	// testMatchUrl := "/en/matches/51ceb4c7/Bastia-Metz-December-10-2016-Ligue-1"
	// testMatchUrl := "/en/matches/7a0a8ab0/Caen-Lille-February-18-2017-Ligue-1 "

	test_6, err_6 := s.DB.GetPlayerRecordsForGivenLeagueMatch(context.Background(), testMatchUrl)
	if err_6 != nil {
		fmt.Println(err_6.Error())
		os.Exit(1)
	}
	fmt.Println("      : Team Name	; 	Player Name 		; 1M ; LM ; Red; Yel; Gol; ")
	for i, item := range test_6 {

		var ind string
		if i < 9 {
			ind = fmt.Sprintf(" %d",i+1)
		} else {
			ind = fmt.Sprintf("%d",i+1)
		}
		tName := fmt.Sprintf(" %s", item.TeamName)
		pName := fmt.Sprintf(" %s %s", item.PlayerName, strings.Repeat(" ",25 - len(item.PlayerName)))

		fmt.Printf("    %s: %s	; %s	; %s ; %s ; %s ; %s ; %s ;\n",
		ind, tName, pName, padStat(item.FirstMinute), padStat(item.LastMinute), padStat(item.RedCard), padStat(item.YellowCard), padStat(item.Goals))

	}
	fmt.Println()
	fmt.Println("Test 6 Complete")
	fmt.Println()
	fmt.Println("===========================================================================================")
}