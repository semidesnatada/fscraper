package analysis

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/semidesnatada/fscraper/config"
	"github.com/semidesnatada/fscraper/database"
)

type PlayerIDPair [2]PlayerID
type SlicePath []PlayerIDPair
type DetailedPath map[PlayerIDPair]RelationshipStats
type RelationshipStats struct {
	CompClubSeasonSummary []RelationshipSeasonSummary
}
type RelationshipSeasonSummary struct {
	GamesPlayedTogether int
	MinutesSharedOnPitch int
	Competition string
	Season string
	Club string
}

func GetPathDetailedStatsFromUrls(s *config.State, urlPath []string) (DetailedPath, error) {

	IDs := IDList{}

	for _, url := range urlPath {
		p2ID, err2 := s.DB.GetPlayerIdFromUrl(context.Background(), url)
		if err2 != nil {
			return DetailedPath{}, err2
		}
		IDs = append(IDs, PlayerID(p2ID))
	}

	out, err := GetPathDetailedStats(s, IDs)//List{PlayerID(p1ID), PlayerID(p2ID)})
	if err != nil {
		return DetailedPath{}, err
	}

	return out, nil
}


func GetPathDetailedStats(s *config.State, path IDList) (DetailedPath, error) {
	formattedPath := SlicePath{}
	var previousID PlayerID
	for i, currentID := range path {
		if i == 0 {
			previousID = currentID
			continue
		}
		newPair := PlayerIDPair{previousID, currentID}
		formattedPath = append(formattedPath, newPair)
		previousID = currentID
	}

	out := make(DetailedPath)
	for _, pair := range formattedPath {
		stats, sErr := GetRelationshipStatsForPlayerIDPair(s, pair)
		if sErr != nil {
			return DetailedPath{}, sErr
		}
		out[pair] = stats
	}

	return out, nil
}

func GetRelationshipStatsForPlayerIDPair(s *config.State, pair PlayerIDPair) (RelationshipStats, error) {

	summaries := []RelationshipSeasonSummary{}

	// do a db query which gets stats for the pair of players

	results, err := s.DB.GetSharedLeagueStatsForTwoPlayersByIDs(context.Background(),
		database.GetSharedLeagueStatsForTwoPlayersByIDsParams{
			ID: uuid.UUID(pair[0]),
			ID_2: uuid.UUID(pair[1]),
		})
	if err != nil {
		return RelationshipStats{}, err
	}

	for _, result := range results {
		sum := RelationshipSeasonSummary {
			GamesPlayedTogether: int(result.SharedMatches),
			MinutesSharedOnPitch: int(result.SharedMinutes),
			Competition: result.CompName,
			Season: result.CompSeason,
			Club: result.TeamName,
		}
		summaries = append(summaries, sum)
	}

	return RelationshipStats{CompClubSeasonSummary: summaries}, nil
}

func (d DetailedPath) PrintPath(s *config.State) error {

	fmt.Println(strings.Repeat("=", 70))

	for pair, stats := range d {

		p1Name, err1 := s.DB.GetPlayerNameFromId(context.Background(), uuid.UUID(pair[0]))
		if err1 != nil {
			return err1
		}
		p2Name, err2 := s.DB.GetPlayerNameFromId(context.Background(), uuid.UUID(pair[1]))
		if err2 != nil {
			return err2
		}
		fmt.Println()
		fmt.Printf("%s and %s played together in the following seasons\n",p1Name, p2Name)

		fmt.Printf("%s; %s; %s; %s; %s\n",
					pSTS("Competition", 20),
					pSTS("Season", 15),
					pSTS("Club", 25),
					pSTS("P",6),
					pSTS("Mins",6))
		for _, stat := range stats.CompClubSeasonSummary {
			fmt.Printf("%s; %s; %s; %s; %s\n",
						pSTS(stat.Competition, 20),
						pSTS(stat.Season, 15),
						pSTS(stat.Club, 25),
						pNTS(stat.GamesPlayedTogether,6),
						pNTS(stat.MinutesSharedOnPitch,6))
		}
		fmt.Println()
	}

	fmt.Println(strings.Repeat("=", 70))
	return nil
}

func pSTS(s string, size int) string {
	//pSTS == padStringToSize
	return fmt.Sprintf("%s%s", s, strings.Repeat(" ", size - len(s)))
}

func pNTS(n int, size int) string {
	//pNTS == padNumberToSize
	var width int
	if n > 99 {
		width = 3
	} else if n > 9 {
		width = 2
	} else {
		width = 1
	}
	return fmt.Sprintf("%d%s", n, strings.Repeat(" ", size - width))
}