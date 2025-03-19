package client

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/semidesnatada/fscraper/config"
	"github.com/semidesnatada/fscraper/database"
)

func storeLeagueMatchSummaries(s *config.State, matches CompetitionSeasonSummary) error {
	compID, compErr := handleComp(s, matches)
	if compErr != nil {
		return compErr
	}

	for _, match := range matches.Data {

		matchParams := database.CreateLeagueMatchParams{
			CompetitionID: compID,
		}

		err := processLeagueMatchSummary(s, match, &matchParams)
		if err != nil {
			if err.Error() == "score not available" {
				continue
				} else {
				return err
			}
		}

		_, dbErr := s.DB.CreateLeagueMatch(
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

func processLeagueMatchSummary(s *config.State, match MatchSummary, temp *database.CreateLeagueMatchParams) (error) {
	goals, gErr, ok := handleScore(match) 
	if gErr != nil {
		return gErr
	}
	if !ok {
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
	temp.Weekday = weekday
	temp.Url = match.data["match-report"]
	temp.HomeTeamOnlineID = match.data["home_team-url"]
	temp.AwayTeamOnlineID = match.data["away_team-url"]

	return  nil
}

func handleTeam(s *config.State, match MatchSummary, homeOrAway bool) (uuid.UUID, error) {
	var teamName string
	var ok bool
	
	if homeOrAway {
		teamName, ok = match.data["home_team"]
	} else {
		teamName, ok = match.data["away_team"]
	}

	if !ok {
		return uuid.UUID{}, errors.New("no team in this record")
	}

	exists, err := s.DB.CheckIfTeamExistsByName(context.Background(), teamName)
	if err != nil {
		return uuid.UUID{}, err
	}

	if exists {
		id, idErr := s.DB.GetTeamIdFromName(
			context.Background(),
			teamName,
		)
		if idErr != nil {
			return uuid.UUID{}, idErr
		}
		return id, nil
	}

	team, teamErr := s.DB.CreateTeam(context.Background(),
	database.CreateTeamParams{
		ID: uuid.New(),
		Name: teamName,
	},
	)
	if teamErr != nil {
		return uuid.UUID{}, teamErr
	}
	return team.ID, nil
}

func handleAttendance(match MatchSummary) (sql.NullInt32, error) {

	attendString, ok := match.data["attendance"]
	if !ok {
		return sql.NullInt32{Valid: false}, nil
	}

	if len(attendString) < 1 {
		return sql.NullInt32{Valid: false}, nil
	}

	attendString = strings.ReplaceAll(attendString, ",", "")

	attendInt, err := strconv.ParseInt(attendString, 10, 32)
	if err != nil {
		return sql.NullInt32{}, err
	}

	newInt := int32(attendInt)

	attendance := sql.NullInt32{
		Int32: newInt,
		Valid: true,
	}

	return attendance, nil
}

func handleXG(match MatchSummary, homeOrAway bool) (sql.NullFloat64, error) {
	var xgString string
	var ok bool

	if homeOrAway{
		xgString, ok = match.data["home_xg"]
	} else {
		xgString, ok = match.data["away_xg"]
	}

	if !ok {
		return sql.NullFloat64{Valid: false}, nil
	}

	if len(xgString) < 1 {
		return sql.NullFloat64{Valid: false}, nil
	}

	float, err := strconv.ParseFloat(xgString, 64)
	if err != nil {
		return sql.NullFloat64{}, err
	}

	xgFloat := sql.NullFloat64{
		Float64: float,
		Valid: true,
	}

	return xgFloat, nil

}

func handleScore(match MatchSummary) ([2]int32, error, bool) {

	scoreString, ok := match.data["score"]

	if !ok {
		return [2]int32{}, errors.New("error accessing score from record"), false
	}
	if len(scoreString) < 1 {
		// fmt.Println("score not found")
		return [2]int32{}, nil, false
	}

	scoreString = strings.ReplaceAll(scoreString, "\u00a0", "")

	goals := strings.Split(scoreString, "–")

	homeGoals, hErr := strconv.Atoi(goals[0])
	if hErr != nil {
		return [2]int32{}, hErr, false
	}
	awayGoals, aErr := strconv.Atoi(goals[1])
	if aErr != nil {
		return [2]int32{}, aErr, false
	}
	return [2]int32{int32(homeGoals), int32(awayGoals)}, nil, true
}

func handleDate(match MatchSummary) (time.Time, error) {

	dateString, ok := match.data["date"]
	if !ok {
		return time.Time{}, errors.New("error accessing date from record")
	}

	date, err := time.Parse(time.DateOnly, dateString)
	if err != nil {
		// fmt.Println(match.data["home_team"], match.data["away_team"], match.data["score"])
		return time.Time{}, err
	}
	return date, nil
}

func handleKickoff(match MatchSummary) (sql.NullTime, error) {
	var kickoff sql.NullTime
	if ko, ok := match.data["start_time"]; ok && len(ko)>0 {
		ko = strings.ReplaceAll(ko, " ", "")
		// kickoffTime, err := time.Parse(time.TimeOnly, ko)
		kickoffTime, err := time.Parse("15:04", ko)
		if err != nil {
			return sql.NullTime{Valid: false}, err
		}
		kickoff = sql.NullTime{
			Time: kickoffTime,
			Valid: true,
		}
	} else {
		kickoff = sql.NullTime{Valid: false}
	}
	return kickoff, nil
}

func handleComp(s *config.State, matches CompetitionSeasonSummary) (uuid.UUID, error) {

	name := matches.CompetitionName
	season := matches.CompetitionSeason
	url := matches.Url

	exists, getErr := s.DB.CheckIfCompetitionExistsByNameAndSeason(
		context.Background(),
		database.CheckIfCompetitionExistsByNameAndSeasonParams{
			Name: name,
			Season: season,
		},
	)
	if getErr != nil {
		return uuid.UUID{}, getErr
	}
	if exists {
		compID, fetchErr := s.DB.GetCompetitionIdFromNameAndSeason(
			context.Background(),
			database.GetCompetitionIdFromNameAndSeasonParams{
				Name: name,
				Season: season,
			},
		)
		if fetchErr != nil {
			return uuid.UUID{}, fetchErr
		}
		return compID, nil
	}

	comp, createErr := s.DB.CreateCompetition(
		context.Background(),
		database.CreateCompetitionParams{
			ID: uuid.New(),
			Name: name,
			Season: season,
			Url: url,
	})
	if createErr != nil {
		return uuid.UUID{}, createErr
	}
	return comp.ID, nil
}


func handleReferee(s *config.State, match MatchSummary) (uuid.NullUUID, error) {
	var referee uuid.NullUUID
	if ref, ok := match.data["referee"]; ok && len(ref) > 0 {

		exists, err := s.DB.CheckIfRefereeExistsByName(context.Background(), ref)
		if err != nil {
			return uuid.NullUUID{Valid: false}, err
		}
		var ref_id uuid.UUID
		if exists {
			ref_id, err = s.DB.GetRefereeIdFromName(context.Background(), ref)
			if err != nil {
				return uuid.NullUUID{Valid: false}, err
			}
		} else {
			ref_record, err := s.DB.CreateReferee(
				context.Background(),
				database.CreateRefereeParams{
					ID: uuid.New(),
					Name: ref,
				})
			if err != nil {
				return uuid.NullUUID{Valid: false},err
			}
			ref_id = ref_record.ID
		}

		referee = uuid.NullUUID{
			UUID: ref_id,
			Valid: true,
		}
	} else {
		referee = uuid.NullUUID{Valid: false}
	}
	return referee, nil
}


func handleVenue(s *config.State, match MatchSummary) (uuid.NullUUID, error) {
	var venue uuid.NullUUID
	if ven, ok := match.data["venue"]; ok && len(ven) > 0 {

		exists, err := s.DB.CheckIfVenueExistsByName(context.Background(), ven)
		if err != nil {
			return uuid.NullUUID{Valid: false}, err
		}
		var ven_id uuid.UUID
		if exists {
			ven_id, err = s.DB.GetVenueIdFromName(context.Background(), ven)
			if err != nil {
				return uuid.NullUUID{Valid: false}, err
			}
		} else {
			ref_record, err := s.DB.CreateVenue(
				context.Background(),
				database.CreateVenueParams{
					ID: uuid.New(),
					Name: ven,
				})
			if err != nil {
				return uuid.NullUUID{Valid: false},err
			}
			ven_id = ref_record.ID
		}

		venue = uuid.NullUUID{
			UUID: ven_id,
			Valid: true,
		}
	} else {
		venue = uuid.NullUUID{Valid: false}
	}
	return venue, nil
}