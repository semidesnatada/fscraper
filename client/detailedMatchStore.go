package client

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/semidesnatada/fscraper/config"
	"github.com/semidesnatada/fscraper/database"
)

func storeDetailedMatches(s *config.State, container DetailedMatchContainer) {

	container = passAndClean(container)

	for _, player := range container.homePlayers {
		template := getRecordTemplate(s, container.matchUrl, container.HomeTeamOnlineID, container.isknockout, player)
		storeTemplate(s, template)
	}
	for _, player := range container.awayPlayers {
		template := getRecordTemplate(s, container.matchUrl, container.AwayTeamOnlineID, container.isknockout, player)
		storeTemplate(s, template)
	}
}

func getRecordTemplate(s *config.State, matchUrl string, teamOnlineID string, isKo bool, player PlayerDetailContainer) database.CreatePlayerMatchParams {

	// checking whether player exists
	player_exists, pErr := s.DB.CheckIfPlayerExistsByUrl(context.Background(), player.player_url)
	if pErr != nil {
		fmt.Println("error checking whether player id exists from the url",pErr.Error())
		os.Exit(1)
	}

	var player_id uuid.UUID

	if player_exists {
		p_id, err := s.DB.GetPlayerIdFromUrl(context.Background(), player.player_url)
		if err != nil {
			fmt.Println("error finding the player id from the url",err.Error())
			os.Exit(1)
		}
		player_id = p_id
	} else {
		playerResponse, err := s.DB.CreatePlayer(context.Background(), database.CreatePlayerParams{
			ID: uuid.New(),
			Name: player.player_name,
			Nationality: player.nationality_link,
			Url: player.player_url,
		})
		if err != nil {
			fmt.Println("error finding the player id from the url",err.Error())
			os.Exit(1)
		}
		player_id = playerResponse.ID
	}

	var isKnockOut bool
	var matchID uuid.UUID
	var mErr error

	if isKo {
		isKnockOut = true
		matchID, mErr = s.DB.GetKnockoutMatchIDFromUrl(context.Background(), matchUrl)
		if mErr != nil {
			fmt.Println("error geting the knockout match id from its url",mErr.Error())
			os.Exit(1)
		}
	} else {
		isKnockOut = false
		matchID, mErr = s.DB.GetLeagueMatchIDFromUrl(context.Background(), matchUrl)
		if mErr != nil {
			fmt.Println("error geting the league match id from its url",mErr.Error())
			os.Exit(1)
		}
	}

	recordTemplate := database.CreatePlayerMatchParams{
		MatchID: matchID,
		PlayerID: player_id,
		MatchUrl: matchUrl,
		FirstMinute: int32(player.first_minute),
		LastMinute: int32(player.last_minute),
		Goals: int32(player.goals),
		Penalties: int32(player.penalties),
		YellowCard: int32(player.yellow_card),
		RedCard: int32(player.red_card),
		OwnGoals: int32(player.own_goals),
		IsKnockout: isKnockOut,
		AtHome: player.home_or_away,
	}

	return recordTemplate
}

func passAndClean(d DetailedMatchContainer) DetailedMatchContainer {

	var last_starter_h int
	var last_starter_last_minute_h int
	last_starter_h = 0
	last_starter_last_minute_h = 0
	//clean minutes data
	for i, player := range d.homePlayers {
		d.homePlayers[i].home_or_away = true

		if player.mins_played == 90 {
			d.homePlayers[i].first_minute = 0
			d.homePlayers[i].last_minute = 90
			last_starter_h = i + 1
			last_starter_last_minute_h = 90
		} else if i == last_starter_h {
			d.homePlayers[i].first_minute = 0
			d.homePlayers[i].last_minute = player.mins_played
			last_starter_last_minute_h = player.mins_played
		} else {
			d.homePlayers[i].first_minute = last_starter_last_minute_h
			d.homePlayers[i].last_minute = last_starter_last_minute_h + player.mins_played
			if d.homePlayers[i].last_minute < 90 {
				last_starter_last_minute_h = d.homePlayers[i].last_minute
			} else {
				last_starter_last_minute_h = 90
				last_starter_h = i + 1
			}
		}
	}
	var last_starter_a int
	var last_starter_last_minute_a int
	last_starter_a = 0
	last_starter_last_minute_a = 0
	for i, player := range d.awayPlayers {
		d.awayPlayers[i].home_or_away = false
		if player.mins_played == 90 {
			d.awayPlayers[i].first_minute = 0
			d.awayPlayers[i].last_minute = 90
			last_starter_a = i + 1
			last_starter_last_minute_a = 90
		} else if i == last_starter_a {
			d.awayPlayers[i].first_minute = 0
			d.awayPlayers[i].last_minute = player.mins_played
			last_starter_last_minute_a = player.mins_played
		} else {
			d.awayPlayers[i].first_minute = last_starter_last_minute_a
			d.awayPlayers[i].last_minute = last_starter_last_minute_a + player.mins_played
			if d.awayPlayers[i].last_minute < 90 {
				last_starter_last_minute_a = d.awayPlayers[i].last_minute
			} else {
				last_starter_last_minute_a = 90
				last_starter_a = i + 1
			}
		}
	}
	return d
}

func storeTemplate(s *config.State, template database.CreatePlayerMatchParams) {

	_, recErr := s.DB.CreatePlayerMatch(
		context.Background(),
		template,
	)
	if recErr != nil {
		fmt.Println("couldn't store prepared record in db",recErr.Error())
		os.Exit(1)
	}

}