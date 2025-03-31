package analysis

import (
	"context"
	"fmt"
	"strings"

	"github.com/semidesnatada/fscraper/config"
)

func TestPlayerMatchData(s *config.State, limit int) {
	handleErr(decorate_test(test_1, 1)(s, limit))
	handleErr(decorate_test(test_2, 2)(s, limit))
	handleErr(decorate_test(test_3, 3)(s, limit*3))
	handleErr(decorate_test(test_4, 4)(s, limit*3))
	handleErr(decorate_test(test_5, 5)(s, limit*3))
	handleErr(decorate_test(test_6, 6)(s, limit))
	handleErr(decorate_test(test_7, 7)(s, limit))
	handleErr(decorate_test(test_8, 8)(s, limit))
}

func padStat(i int32) string {

	if i < 10 {
		return fmt.Sprintf(" %d",i)
	} else {
		return fmt.Sprintf("%d", i)
	}

}

func handleErr(e error) {
	if e != nil {
		fmt.Println(e.Error())
	}
}

func decorate_test(f func(s *config.State, l int) error, no int) func(*config.State, int) error {

	return func(s *config.State,l int) error {
		fmt.Println("===========================================================================================")
		fmt.Println()
		fmt.Printf("Beginning Test %d\n", no)
		fmt.Println()
		err := f(s, l)
		if err != nil {
			fmt.Println()
			fmt.Printf("Test %d Failed\n", no)
			fmt.Println()
			fmt.Println("===========================================================================================")
		} else {
			fmt.Println()
			fmt.Printf("Test %d Complete\n", no)
			fmt.Println()
			fmt.Println("===========================================================================================")
		}
		return err
	}
}

func test_1(s *config.State, limit int) error {
	test, err := s.DB.GetNumberOfDistinctPlayersFieldedInLeagueByTeam(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("      : Competition		; 	Team				;Matches Played	;Players Fielded	")
	for i, item := range test {
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
	return nil
}

func test_2(s *config.State, limit int) error {
	test, err := s.DB.GetNumberOfGoalsScoredInEachLeagueSeasonByTeam(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("      : Competition			; 	Team			;Matches Played	;Goals Scored	;Players Fielded	")
	for i, item := range test {
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
	return nil
}

func test_3(s *config.State, limit int) error {
	test, err := s.DB.GetLeagueAllTimeTopScorers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("      : Player			; 	Competition			;Matches Played	;Goals Scored	")
	for i, item := range test {
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
	return nil
}

func test_4(s *config.State, limit int) error {
	test, err := s.DB.GetPlayerLeagueStats(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("      : Player			; 	Competition			;Matches Played	;Goals Scored	")
	for i, item := range test {
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
	return nil
}

func test_5(s *config.State, limit int) error {

	test, err := s.DB.GetAllTimeTopScorers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("      : Player				;Matches Played	;Goals Scored	")
	for i, item := range test {
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
	return nil
}

func test_6(s *config.State, limit int) error {
	// testUrl := "/en/players/c596fcb0/Callum-Wilson"
	testUrl := "/en/players/d70ce98e/Lionel-Messi"
	// testUrl := "/en/players/2c6835e5/Lewis-Miley"
	// testUrl := "/en/players/2b81295d/Raul"

	test_pre, err_pre := s.DB.GetStatsForPlayerUrl(context.Background(), testUrl)
	if err_pre != nil {
		return err_pre
	}
	test, err := s.DB.GetPlayersPlayedWithByUrl(context.Background(), testUrl)
	if err != nil {
		return err
	}
	fmt.Println()
	fmt.Println("Player stats are:")
	fmt.Printf("Name: %s , calculated mins played: %d\n", test_pre.PlayerName, test[0].TotalMinsPlayed)

	fmt.Println()
	fmt.Println("Stats calculated with independent method:")
	fmt.Println("Matches played :", test_pre.MatchesPlayed)
	fmt.Println("Total mins played :", test_pre.TotalMinsPlayed)
	fmt.Println("Goals scorewd :", test_pre.TotalGoals)
	fmt.Println("Pens scored :", test_pre.TotalPens)
	fmt.Println("Own goals scored :", test_pre.TotalOgs)
	fmt.Println("Total yellow cards :", test_pre.TotalYellowCard)
	fmt.Println("Total red cards :", test_pre.TotalRedCard)
	
	fmt.Println()
	fmt.Println("      : Colleague 				; Mins played together")
	for i, item := range test {
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
	return nil
}

func test_7(s *config.State, limit int) error {
	
	test, err := s.DB.GetMatchesWhereMinsDontAddUp(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("      : Total Mins			; 	Match Url")
	for i, item := range test {
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
	return nil
}

func test_8(s *config.State, _ int) error {

	// testMatchUrl := "/en/matches/75e1f5c4/Valencia-Mallorca-November-28-2004-La-Liga"
	// testMatchUrl := "/en/matches/198facdd/Stuttgart-Wolfsburg-September-16-2017-Bundesliga"
	// testMatchUrl := "/en/matches/64b77eb3/Barcelona-Deportivo-La-Coruna-March-18-2000-La-Liga"
	// testMatchUrl := "/en/matches/474ce5b7/Schalke-04-Bayer-Leverkusen-April-9-1996-Bundesliga"
	// testMatchUrl := "/en/matches/51ceb4c7/Bastia-Metz-December-10-2016-Ligue-1"
	// testMatchUrl := "/en/matches/7a0a8ab0/Caen-Lille-February-18-2017-Ligue-1 "
	testMatchUrl := "/en/matches/27e1ab19/Monaco-Lyon-May-2-2021-Ligue-1"

	test, err := s.DB.GetPlayerRecordsForGivenLeagueMatch(context.Background(), testMatchUrl)
	if err != nil {
		return err
	}
	fmt.Println("      : Team Name	; 	Player Name 		; 1M ; LM ; Red; Yel; Gol; ")
	for i, item := range test {

		var ind string
		if i < 9 {
			ind = fmt.Sprintf(" %d",i+1)
		} else {
			ind = fmt.Sprintf("%d",i+1)
		}
		tName := fmt.Sprintf(" %s", item.TeamName)
		pName := fmt.Sprintf(" %s %s", item.PlayerName, strings.Repeat(" ",30 - len(item.PlayerName)))

		fmt.Printf("    %s: %s	; %s	; %s ; %s ; %s ; %s ; %s ;\n",
		ind, tName, pName, padStat(item.FirstMinute), padStat(item.LastMinute), padStat(item.RedCard), padStat(item.YellowCard), padStat(item.Goals))

	}
	return nil
}