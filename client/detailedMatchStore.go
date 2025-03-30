package client

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/google/uuid"
	"github.com/semidesnatada/fscraper/config"
	"github.com/semidesnatada/fscraper/database"
)

func storeDetailedMatches(s *config.State, container DetailedMatchContainer) {

	container = passAndClean(container)

	for _, player := range container.homePlayers {
		template, tempErr := getRecordTemplate(s, container.matchUrl, container.HomeTeamOnlineID, container.isknockout, player)
		if tempErr == nil {
			if template.RedCard == 1 {
				//create and store another template indicating a red card 
				place := createRedCardPlaceholderTemplate(template)
				_, err := createRedCardPlayerRecord(s, place.PlayerID)
				if err != nil {
					fmt.Printf("issue in resolving a red card in the database for match %s, %v", template.MatchUrl, err.Error())
				}
				storeTemplate(s, place)
			}
			storeTemplate(s, template)
		} else {fmt.Printf("error in getting a template for player %s, %v\n", player.player_name, tempErr.Error())}
	}
	for _, player := range container.awayPlayers {
		template, tempErr := getRecordTemplate(s, container.matchUrl, container.AwayTeamOnlineID, container.isknockout, player)
		if tempErr == nil {
			if template.RedCard == 1 {
				//create and store another template indicating a red card 
				place := createRedCardPlaceholderTemplate(template)
				_, err := createRedCardPlayerRecord(s, place.PlayerID)
				if err != nil {
					fmt.Printf("issue in resolving a red card in the database for match %s, %v", template.MatchUrl, err.Error())
				}
				storeTemplate(s, place)
			}
			storeTemplate(s, template)
		} else {fmt.Printf("error in getting a template for player %s, %v\n", player.player_name, tempErr.Error())}
	}
	fmt.Printf("==== Completed storing url: %s \n", container.matchUrl)
	fmt.Println("__________")
}

func getRecordTemplate(s *config.State, matchUrl string, teamOnlineID string, isKo bool, player PlayerDetailContainer) (database.CreatePlayerMatchParams, error) {

	// checking whether player exists
	player_exists, pErr := s.DB.CheckIfPlayerExistsByUrl(context.Background(), player.player_url)
	if pErr != nil {
		return database.CreatePlayerMatchParams{}, fmt.Errorf("error checking whether player id (for %s) exists from the url %w", player.player_name, pErr)
	}

	// get the right player from the database, or create the player if they don't exist
	var player_id uuid.UUID
	if player_exists {
		p_id, err := s.DB.GetPlayerIdFromUrl(context.Background(), player.player_url)
		if err != nil {
			return database.CreatePlayerMatchParams{}, fmt.Errorf("error finding the player id (for %s) from the url %w", player.player_name, err)
		}
		player_id = p_id
	} else {
		playerResponse, cErr := createPlayerRecord(s, player)
		if cErr != nil {
			return database.CreatePlayerMatchParams{}, cErr
		}
		player_id = playerResponse.ID
	}

	matchID, mErr := getMatchID(s, matchUrl, isKo)
	if mErr != nil {
		return database.CreatePlayerMatchParams{}, fmt.Errorf("error geting the match id from its url: %w", mErr)
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
		IsKnockout: isKo,
		AtHome: player.home_or_away,
	}
	return recordTemplate, nil
}

func createRedCardPlaceholderTemplate(parent database.CreatePlayerMatchParams) database.CreatePlayerMatchParams {

	return database.CreatePlayerMatchParams{
		MatchID: parent.MatchID,
		PlayerID: uuid.New(),
		MatchUrl: parent.MatchUrl,
		FirstMinute: parent.LastMinute,
		LastMinute: int32(90),
		Goals: 0,
		Penalties: 0,
		YellowCard: 0,
		RedCard: 0,
		OwnGoals: 0,
		IsKnockout: parent.IsKnockout,
		AtHome: parent.AtHome,
	}
}

func getMatchID(s *config.State, matchUrl string, isKo bool) (uuid.UUID, error)  {
	if isKo {
		matchID, mErr := s.DB.GetKnockoutMatchIDFromUrl(context.Background(), matchUrl)
		return matchID, mErr

	} else {
		matchID, mErr := s.DB.GetLeagueMatchIDFromUrl(context.Background(), matchUrl)
		return matchID, mErr
	}
}

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func createRedCardPlayerRecord(s *config.State, id uuid.UUID) (database.Player, error) {
	// creates a fake player record to occupy the minutes for a player who receives a red card
	placeholderPlayer, pErr := s.DB.CreatePlayer(context.Background(), database.CreatePlayerParams{
		ID: id,
		Name: "fakeRedCard",
		Nationality: "sharmika",
		Url: "milffitness.com"+randSeq(5),
	})
	if pErr == nil {
		return placeholderPlayer, nil
	} else {
		return database.Player{}, fmt.Errorf("error adding a fake red card placeholder in the db: %w", pErr)
	}
} 

func createPlayerRecord(s *config.State, player PlayerDetailContainer) (database.Player, error) {
	playerResponse, err := s.DB.CreatePlayer(context.Background(), database.CreatePlayerParams{
		ID: uuid.New(),
		Name: player.player_name,
		Nationality: player.nationality_link,
		Url: player.player_url,
	})
	if err == nil {
		return playerResponse, nil
	}

	//instead of returning the error here, we use a wildcard player
	//return database.CreatePlayerMatchParams{}, fmt.Errorf("error a new player (%s) in the db: %w", player.player_name, err)
	wildcardPlayer, wErr := s.DB.CreatePlayer(context.Background(), database.CreatePlayerParams{
		ID: uuid.New(),
		Name: "wildcard",
		Nationality: "kuwandrandian",
		Url: "nationalgeographic.com"+randSeq(5),
	})
	if wErr == nil {
		return wildcardPlayer, nil
	} else {
		return database.Player{}, fmt.Errorf("error adding a new player (%s) in the db: %w", player.player_name, wErr)
	}

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
			// if a player plays 90 mins then we automatically know their start and end minute
			d.homePlayers[i].first_minute = 0
			d.homePlayers[i].last_minute = 90
			last_starter_h = i + 1
			last_starter_last_minute_h = 90
		} else if i == last_starter_h {
			// if the helper int tells us this player is a starter, then the start and end minutes are also known
			d.homePlayers[i].first_minute = 0
			d.homePlayers[i].last_minute = player.mins_played
			last_starter_last_minute_h = player.mins_played
			if player.red_card == 1 {
				// if the player received a red card, the next player in the list must have started (they can't be a sub for a red carded player)
				last_starter_h = i + 1
				last_starter_last_minute_h = 90
			}
		} else {
			// otherwise, the player must have been brough on as a sub
			if last_starter_last_minute_h + player.mins_played > 90 {
				// if this holds true, the current player did not in fact come on as a sub, and so the last player went off for a reason other than a red card
				// this means the current player must have started, and the next player will come on as a sub
				d.homePlayers[i].first_minute = 0
				d.homePlayers[i].last_minute = player.mins_played
				last_starter_last_minute_h = player.mins_played
			} else {
			d.homePlayers[i].first_minute = last_starter_last_minute_h
			d.homePlayers[i].last_minute = last_starter_last_minute_h + player.mins_played
			if d.homePlayers[i].last_minute < 90 {
				// this means either the player was sent off, was subbed off again, or didn't play for some other reason
				if player.red_card == 1 {
					last_starter_h = i + 1
					last_starter_last_minute_h = 90
				} else {
				last_starter_last_minute_h = d.homePlayers[i].last_minute
				}
			} else {
				last_starter_last_minute_h = 90
				last_starter_h = i + 1
			}
		}}
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
			if player.red_card == 1 {
				// if the player received a red card, the next player in the list must have started (they can't be a sub for a red carded player)
				last_starter_a = i + 1
				last_starter_last_minute_a = 90
			}
		} else {
			// otherwise, the player must have been brough on as a sub
			if last_starter_last_minute_a + player.mins_played > 90 {
				// if this holds true, the current player did not in fact come on as a sub, and so the last player went off for a reason other than a red card
				// this means the current player must have started, and the next player will come on as a sub
				d.awayPlayers[i].first_minute = 0
				d.awayPlayers[i].last_minute = player.mins_played
				last_starter_last_minute_a = player.mins_played
			} else {
			d.awayPlayers[i].first_minute = last_starter_last_minute_a
			d.awayPlayers[i].last_minute = last_starter_last_minute_a + player.mins_played
			if d.awayPlayers[i].last_minute < 90 {
				// this means either the player was sent off, was subbed off again, or didn't play for some other reason
				if player.red_card == 1 {
					last_starter_a = i + 1
					last_starter_last_minute_a = 90
				} else {
					last_starter_last_minute_a = d.awayPlayers[i].last_minute
				}
			} else {
				last_starter_last_minute_a = 90
				last_starter_a = i + 1
			}
		}}
	}
	return d
}

func storeTemplate(s *config.State, template database.CreatePlayerMatchParams) {

	_, recErr := s.DB.CreatePlayerMatch(
		context.Background(),
		template,
	)
	if recErr != nil {
		fmt.Printf("couldn't store prepared record in db: %s, %s\n", template.MatchUrl, recErr.Error())
		// write the template to json
		bytes, err1 := json.Marshal(template)
		if err1 != nil {
			fmt.Println("failed to marshal erroneous record to json")
		}
		dir, _ := os.Getwd()
		err2 := os.WriteFile(dir+fmt.Sprintf("/db_player_matches_write_errors/id_%s.json", template.MatchID), bytes, 0666)
		if err2 != nil {
			fmt.Println("failed to write erroneous record to json file")
			fmt.Println(err2.Error())
		}
		// os.Exit(1)
	}

}