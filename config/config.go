package config

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/semidesnatada/fscraper/database"
)


type State struct {
	DB *database.Queries
	DBURL string
}

func CreateState(DB_URL string) State {

	db, err := sql.Open("postgres", DB_URL)
	if err != nil {
		fmt.Printf("Programme terminated: %s\n", err.Error())
		os.Exit(1)
	}
	dbQueries := database.New(db)

	appState := State{
		DBURL: "postgres://seanlowery:@localhost:5432/fscraped?sslmode=disable",
		DB: dbQueries,
	}

	return appState
}

func (s *State) DeleteSummaryDatabases() error {

	errorComp := s.DB.DeleteCompetitions(context.Background())
	if errorComp != nil {
		return errorComp
	}
	errorRefs := s.DB.DeleteReferees(context.Background())
	if errorRefs != nil {
		return errorRefs
	}
	errorLMatch := s.DB.DeleteLeagueMatches(context.Background())
	if errorLMatch != nil {
		return errorLMatch
	}
	errorKMatch := s.DB.DeleteKnockoutMatches(context.Background())
	if errorKMatch != nil {
		return errorKMatch
	}
	errorTeam := s.DB.DeleteTeams(context.Background())
	if errorTeam != nil {
		return errorTeam
	}
	errorVen := s.DB.DeleteVenues(context.Background())
	if errorVen != nil {
		return errorVen
	}
	return nil
}

func (s *State) DeleteDetailedDatabases() error {

	errorPm := s.DB.DeletePlayerMatches(context.Background())
	if errorPm != nil {
		return errorPm
	}
	errorP := s.DB.DeletePlayers(context.Background())
	if errorP != nil {
		return errorP
	}

	return nil
}

func (s *State) DeleteAllDatabases() error {
	dErr := s.DeleteDetailedDatabases()
	if dErr != nil {
		return dErr
	}
	sErr := s.DeleteSummaryDatabases()
	if sErr != nil {
		return sErr
	}
	return nil
}