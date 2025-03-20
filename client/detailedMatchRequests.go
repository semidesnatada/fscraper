package client

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/semidesnatada/fscraper/config"
	"github.com/semidesnatada/fscraper/database"
)

func DetailedMatchScraper(s *config.State) error {

	// get match urls via query to database
	// leagueMatchUrls, lErr := s.DB.GetLeagueMatchUrlsAndTeamOnlineIds(context.Background())
	// if lErr != nil {
	// 	return lErr
	// }
	knockoutMatchUrls, kErr := s.DB.GetKnockoutMatchUrlsAndTeamOnlineIds(context.Background())
	if kErr != nil {
		return kErr
	}

	// scrapeLeaguesMatches(s, leagueMatchUrls)
	scrapeKnockoutsMatches(s, knockoutMatchUrls)

	return nil

}

// func scrapeLeaguesMatches(s *config.State, leagueMatchUrls []database.GetLeagueMatchUrlsAndTeamOnlineIdsRow) error {

// 	return nil

// }

func scrapeKnockoutsMatches(s *config.State, knockoutMatchUrls []database.GetKnockoutMatchUrlsAndTeamOnlineIdsRow) error {

	processedMatchChannel := make(chan DetailedMatchContainer)

	// create a go routine to scrape all the urls and allow storing in the db concurrently.
	go processKnockoutUrlBlock(s, knockoutMatchUrls, processedMatchChannel)

	//add the record to the database.
	for processed := range processedMatchChannel {
		//store matchcontainer in database here
		fmt.Printf("=== Started storing url: %s \n", processed.matchUrl)

		storeDetailedMatches(s, processed)
		
		fmt.Printf("==== Completed storing url: %s \n", processed.matchUrl)

	}
	return nil
}

func processKnockoutUrlBlock(s *config.State, urls []database.GetKnockoutMatchUrlsAndTeamOnlineIdsRow, c chan DetailedMatchContainer) {
	ticker := time.NewTicker(time.Second * 6)
	// iterate over each url in turn
	for _, row := range urls {

		fmt.Printf("= Started scraping url: %s \n", row.Url)

		//check if url is already present in db
		exists, eErr := s.DB.CheckIfMatchUrlExistsInPlayerMatches(context.Background(),row.Url)
		if eErr != nil {
			fmt.Println(eErr.Error())
			os.Exit(1)
		}
		//if it is already present, skip to next iteration
		if exists {
			fmt.Println("== Already scraped, moving on")
			continue
		}

		// get response from url
		combinedUrl := baseMatchUrl+row.Url
		res, resErr := http.Get(combinedUrl)
		if resErr != nil {
			fmt.Println(resErr.Error())
			os.Exit(1)
		}

		//parse the response with a goroutine. parsing should take much less 
		// than six seconds, so should complete before the ticker is reached anyways.
		go parseKnockoutResponse(res, c, row)
		//block until enough time has passed, to rate limit http requests.
		<- ticker.C
	}
	close(c)
}

func parseKnockoutResponse(res *http.Response, c chan DetailedMatchContainer, row database.GetKnockoutMatchUrlsAndTeamOnlineIdsRow) {

	homeCode := strings.Split(row.HomeTeamOnlineID,"/")[3]
	awayCode := strings.Split(row.AwayTeamOnlineID,"/")[3]

	var homePlayers []PlayerDetailContainer
	var awayPlayers []PlayerDetailContainer

	defer res.Body.Close()
	processDocAndTeam(res, &homePlayers, &awayPlayers, homeCode, awayCode)
	// processDocAndTeam(res, &awayPlayers, awayCode)

	processedResult := DetailedMatchContainer{
		matchUrl: row.Url,
		HomeTeamOnlineID: homeCode,
		AwayTeamOnlineID: awayCode,
		homePlayers: homePlayers,
		awayPlayers: awayPlayers,
		isknockout: true,
	}

	fmt.Printf("== Finished scraping url: %s successfully\n", row.Url)
	c <- processedResult
}
