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

func StoreMatchSummaries(s *config.State, matches CompetitionSeasonSummary) error {

	for _, match := range matches.Data {
		// compID, compErr := s.DB.GetCompetitionIdFromName(
		// 	context.Background(),
		// 	matches.CompetitionName,
		// )
		// if compErr != nil {
		// 	return compErr
		// }

		matchParams := database.CreateMatchParams{}

		err := processMatch(s, match, &matchParams)
		if err != nil {
			return err
		}

		// fmt.Println(matchParams)

		_, dbErr := s.DB.CreateMatch(
			context.Background(),
			matchParams,
		)

		// fmt.Println(matchino)
		if dbErr != nil {
			fmt.Println(dbErr.Error())
			return dbErr
		}

	}
	return nil
}

func processMatch(s *config.State, match MatchSummary, temp *database.CreateMatchParams) ( error) {
	compID := uuid.New()
	homeID, homeErr :=  handleHome(s, match, true)
	if homeErr != nil {
		return homeErr
	}
	// fmt.Println("success: home_team")
	awayID, awayErr :=  handleHome(s, match, false)
	if awayErr != nil {
		return awayErr
	}
	// fmt.Println("success: away_team")
	date, dateErr := handleDate(match)
	if dateErr != nil {
		return dateErr
	}
	// fmt.Println("success: date")
	goals, gErr := handleScore(match) 
	if gErr != nil {
		return gErr
	}
	// fmt.Println("success: score")
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

	temp.ID = uuid.New()
	temp.CompetitionID = compID
	temp.CompetitionSeasonID = "21-22"
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

	return  nil
}

func handleHome(s *config.State, match MatchSummary, homeOrAway bool) (uuid.UUID, error) {
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

func handleScore(match MatchSummary) ([2]int32, error) {

	scoreString, ok := match.data["score"]

	if !ok {
		return [2]int32{}, errors.New("error accessing score from record")
	}

	goals := strings.Split(scoreString, "–")

	homeGoals, hErr := strconv.Atoi(goals[0])
	if hErr != nil {
		return [2]int32{}, hErr
	}
	awayGoals, aErr := strconv.Atoi(goals[1])
	if aErr != nil {
		return [2]int32{}, aErr
	}
	return [2]int32{int32(homeGoals), int32(awayGoals)}, nil
}

func handleDate(match MatchSummary) (time.Time, error) {

	dateString, ok := match.data["date"]
	if !ok {
		return time.Time{}, errors.New("error accessing date from record")
	}

	date, err := time.Parse(time.DateOnly, dateString)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}

func handleKickoff(match MatchSummary) (sql.NullTime, error) {
	var kickoff sql.NullTime
	if ko, ok := match.data["start_time"]; ok {
		kickoffTime, err := time.Parse(time.TimeOnly, ko)
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

func handleReferee(s *config.State, match MatchSummary) (uuid.NullUUID, error) {
	var referee uuid.NullUUID
	if ref, ok := match.data["referee"]; ok {

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
	if ven, ok := match.data["venue"]; ok {

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