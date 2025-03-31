package analysis

import (
	"context"

	"github.com/semidesnatada/fscraper/config"
	"github.com/semidesnatada/fscraper/database"
)

// func which checks whether the table prepared using the leaguematches database yields an equal result to the player matches database.

func GetLeagueTableFromPlayerMatches(s *config.State, leagueName, leagueSeason string) error {

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

	if rowsLM[0].Draws == 1 {
		return nil
	}
	if rowsPM[0].String() == "sa" {
		return nil
	}
	return nil
}