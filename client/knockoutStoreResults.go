package client

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/semidesnatada/fscraper/config"
	"github.com/semidesnatada/fscraper/database"
)

func storeKnockoutMatchSummaries(s *config.State, matches CompetitionSeasonSummary) error {
	compID, compErr := handleComp(s, matches)
	if compErr != nil {
		return compErr
	}

	for _, match := range matches.Data {

		matchParams := database.CreateKnockoutMatchParams{
			CompetitionID: compID,
		}

		err := processKnockoutMatchSummary(s, match, &matchParams)
		if err != nil {
			return err
		}

		_, dbErr := s.DB.CreateKnockoutMatch(
			context.Background(),
			matchParams,
		)

		if dbErr != nil {
			fmt.Println(dbErr.Error())
			return dbErr
		}
	}
	return nil
}

func processKnockoutMatchSummary(s *config.State, match MatchSummary, temp *database.CreateKnockoutMatchParams) (error) {
	goals, goalsOk, pens, wasShootout, gErr := handleKnockoutScore(match) 
	// fmt.Println()
	// fmt.Println("==========")
	// fmt.Println(goals)
	// fmt.Println(goalsOk)
	// fmt.Println(pens)
	// fmt.Println(wasShootout)
	// fmt.Println(match.data["home_team"])
	// fmt.Println(match.data["away_team"])
	// fmt.Println(match)
	if gErr != nil {
		return gErr
	}
	if !goalsOk {

		return errors.New("score not available")
	}
	// fmt.Println("success: score")
	homeID, homeErr :=  handleTeam(s, match, true)
	if homeErr != nil {
		return homeErr
	}
	// fmt.Printf("success: home_team %s \n", match.data["home_team"])
	awayID, awayErr :=  handleTeam(s, match, false)
	if awayErr != nil {
		return awayErr
	}
	// fmt.Printf("success: away_team %s \n", match.data["away_team"])
	date, dateErr := handleDate(match)
	if dateErr != nil {
		return dateErr
	}
	// fmt.Println("success: date")
	kickoff, koErr := handleKickoff(match)
	if koErr != nil {
		return koErr
	}
	// fmt.Println("success: ko_time")
	referee, refErr := handleReferee(s, match)
	if refErr != nil {
		return refErr
	}
	// fmt.Println("success: referee")
	venue, venErr := handleVenue(s, match)
	if venErr != nil {
		return venErr
	}
	// fmt.Println("success: venue")
	attendance, attendErr := handleAttendance(match)
	if attendErr != nil {
		return attendErr
	}
	// fmt.Println("success: attendance")
	homeXG, hxgErr := handleXG(match, true)
	if hxgErr != nil {
		return hxgErr
	} 
	// fmt.Println("success: home_xg")
	awayXG, axgErr := handleXG(match, false)
	if axgErr != nil {
		return axgErr
	} 
	// fmt.Println("success: away_xg")
	weekday, ok := match.data["dayofweek"]
	if !ok {
		return errors.New("error accessing day of week")
	}
	// fmt.Println("success: weekday")

	// fmt.Println(match.data["match-report"])

	round, roundErr := handleRound(match)
	if roundErr != nil {
		return roundErr
	}

	temp.ID = uuid.New()
	temp.HomeTeamID = homeID
	temp.AwayTeamID = awayID
	temp.HomeGoals = goals[0]
	temp.AwayGoals = goals[1]
	temp.Date = date
	temp.KickOffTime = kickoff
	temp.RefereeID = referee
	temp.VenueID = venue
	temp.Attendance = attendance
	temp.HomeXg = homeXG
	temp.AwayXg = awayXG
	temp.WentToPens = wasShootout
	temp.HomeGoals = goals[0]
	temp.AwayGoals = goals[1]
	if wasShootout {
		temp.HomePens = sql.NullInt32{Int32:pens[0], Valid:true}
		temp.AwayPens = sql.NullInt32{Int32:pens[1], Valid:true}
	} else {
		temp.HomePens = sql.NullInt32{Valid:false}
		temp.AwayPens = sql.NullInt32{Valid:false}
	}
	temp.Round = round
	temp.Weekday = weekday
	temp.Url = match.data["match-report"]

	return  nil
}

func handleKnockoutScore(m MatchSummary) ([2]int32, bool, [2]int32, bool, error) {

	scoreString, ok := m.data["score"]
	if !ok {
		return [2]int32{}, false, [2]int32{}, false, errors.New("errors accessing score from record")
	}
	if len(scoreString) < 1 {
		// fmt.Println("score not found")
		return [2]int32{}, false, [2]int32{}, false, nil
	}
	scoreString = strings.ReplaceAll(scoreString, "\u00a0", "")

	var goals [2]int32
	var pens [2]int32

	if pOk := strings.Contains(scoreString, "(") && strings.Contains(scoreString, ")"); pOk {
		goalPiece := strings.Split(scoreString, "–")

		v1, err1 := strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(strings.Fields(goalPiece[0])[0], "(", ""), ")", ""))
		v2, err2 := strconv.Atoi(strings.Fields(goalPiece[0])[1])
		v3, err3 := strconv.Atoi(strings.Fields(goalPiece[1])[0])
		v4, err4 := strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(strings.Fields(goalPiece[1])[1], "(", ""), ")", ""))

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			fmt.Printf("error in %s vs %s\n", m.data["home_team"], m.data["away_team"])
			return [2]int32{}, false, [2]int32{}, false, errors.New("error parsing the score from a knockout match")
		}
		goals[0] = int32(v2)
		goals[1] = int32(v3)
		pens[0] = int32(v1)
		pens[1] = int32(v4)
		return goals, ok, pens, pOk, nil 
	} else {
		goalPieces := strings.Split(scoreString, "–")

		homeGoals, hErr := strconv.Atoi(goalPieces[0])
		if hErr != nil {
			return [2]int32{}, false, [2]int32{}, false, hErr
		}
		awayGoals, aErr := strconv.Atoi(goalPieces[1])
		if aErr != nil {
			return [2]int32{}, false,[2]int32{}, false, aErr
		}
		return [2]int32{int32(homeGoals), int32(awayGoals)}, ok, [2]int32{}, false, nil
	}
}

func handleRound(m MatchSummary) (string, error) {
	round, ok := m.data["round"]
	if !ok {
		return "", errors.New("no round data provided")
	}
	return round, nil
}