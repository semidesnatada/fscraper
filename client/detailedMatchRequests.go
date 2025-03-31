package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/semidesnatada/fscraper/config"
	"github.com/semidesnatada/fscraper/database"
)

func DetailedMatchScraper(s *config.State) error {

	// get match urls via query to database
	// leagueMatchUrls, lErr := s.DB.GetLeagueMatchUrlsAndTeamOnlineIdsWOffset(context.Background(), 0)
	leagueMatchUrls_f, lErr := s.DB.GetLeagueMatchUrlsAndTeamOnlineIds(context.Background())
	if lErr != nil {
		return lErr
	}
	// knockoutMatchUrls, kErr := s.DB.GetKnockoutMatchUrlsAndTeamOnlineIdsWOffset(context.Background(), 0)
	// if kErr != nil {
	// 	return kErr
	// }

	leagueMatchUrls := []database.GetLeagueMatchUrlsAndTeamOnlineIdsWOffsetRow{}
	for _, m := range leagueMatchUrls_f {
		leagueMatchUrls = append(leagueMatchUrls, database.GetLeagueMatchUrlsAndTeamOnlineIdsWOffsetRow{Url:m.Url, HomeTeamOnlineID: m.HomeTeamOnlineID , AwayTeamOnlineID: m.AwayTeamOnlineID})
	}

	leErr := scrapeLeaguesMatches(s, leagueMatchUrls)
	if leErr != nil {
		return leErr
	}
	// koErr := scrapeKnockoutsMatches(s, knockoutMatchUrls)
	// if koErr != nil {
	// 	return koErr
	// }

	return nil

}

//top of listt
// /en/matches/00023357/Hoffenheim-Bayer-Leverkusen-April-12-2021-Bundesliga
//  /en/matches/0005cd5f/Napoli-Juventus-December-1-2017-Serie-A
//  /en/matches/0006415c/Udinese-Atalanta-October-9-2022-Serie-A

// 1000 and 1001
// /en/matches/0454c687/Levante-Espanyol-April-15-2016-La-Liga
//  /en/matches/0454ebc4/Albacete-Lleida-March-27-1994-La-Liga

func scrapeLeaguesMatches(s *config.State, leagueMatchUrls []database.GetLeagueMatchUrlsAndTeamOnlineIdsWOffsetRow) error {

	processedMatchChannel := make(chan DetailedMatchContainer)

	// create a go routine to scrape all the urls and allow storing in the db concurrently.
	go processLeagueUrlBlock(s, leagueMatchUrls, processedMatchChannel)

	//add the record to the database.
	for processed := range processedMatchChannel {
		//store matchcontainer in database here
		fmt.Printf("=== Started storing url: %s \n", processed.matchUrl)

		go storeDetailedMatches(s, processed)
		
	}
	return nil
}

func scrapeKnockoutsMatches(s *config.State, knockoutMatchUrls []database.GetKnockoutMatchUrlsAndTeamOnlineIdsWOffsetRow) error {

	processedMatchChannel := make(chan DetailedMatchContainer)

	// create a go routine to scrape all the urls and allow storing in the db concurrently.
	go processKnockoutUrlBlock(s, knockoutMatchUrls, processedMatchChannel)

	//add the record to the database.
	for processed := range processedMatchChannel {
		//store matchcontainer in database here
		fmt.Printf("=== Started storing url: %s \n", processed.matchUrl)

		go storeDetailedMatches(s, processed)
		
		fmt.Printf("==== Completed storing url: %s \n", processed.matchUrl)
		fmt.Println("__________")

	}
	return nil
}

func processLeagueUrlBlock(s *config.State, urls []database.GetLeagueMatchUrlsAndTeamOnlineIdsWOffsetRow, c chan DetailedMatchContainer) {
	ticker := time.NewTicker(time.Second * 6)

	// ideally we should do one sql query up here and check against the result in the loop below,
	// rather than one query for each run of the loop.
	// this would significantly speed up the checking loops, where most matches are already 
	// present in the db, and the rate limiting is skipped.

	// iterate over each url in turn
	for i, row := range urls {
		
		fmt.Printf("= %d Started scraping url: %s \n", i, row.Url)

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

		// write response to disk in order to save time from rate limiting in future
		dir, _ := os.Getwd()
		// fmt.Println(dir+"/db_webpage_store/"+strings.Split(row.Url,"/")[3]+".html")
		filePath := dir+"/db_webpage_store/"+strings.Split(row.Url,"/")[3]+".html"
		fileExists := checkFileExists(filePath)
		if fileExists {
			// do nothing for now, but should simply skip the http request in future
		} else {
			// f, err := os.Create(filePath)
			// if err != nil {
			// 	fmt.Println(err.Error())
			// 	os.Exit(1)
			// }
			// defer f.Close()
			dump, _ := httputil.DumpResponse(res, true)
			err2 := os.WriteFile(filePath, dump, 0644)
			// err2 := res.Write(f)
			if err2 != nil {
				fmt.Println(err2.Error())
				os.Exit(1)
			}
		}
		//parse the response with a goroutine. parsing should take much less 
		// than six seconds, so should complete before the ticker is reached anyways.
		go parseLeagueResponse(res, c, row)
		//block until enough time has passed, to rate limit http requests.
		<- ticker.C
	}
	close(c)
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}

func processKnockoutUrlBlock(s *config.State, urls []database.GetKnockoutMatchUrlsAndTeamOnlineIdsWOffsetRow, c chan DetailedMatchContainer) {
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

func parseKnockoutResponse(res *http.Response, c chan DetailedMatchContainer, row database.GetKnockoutMatchUrlsAndTeamOnlineIdsWOffsetRow) {

	homeCode := strings.Split(row.HomeTeamOnlineID,"/")[3]
	awayCode := strings.Split(row.AwayTeamOnlineID,"/")[3]

	var homePlayers []PlayerDetailContainer
	var awayPlayers []PlayerDetailContainer

	defer res.Body.Close()
	processDocAndTeam(res, &homePlayers, &awayPlayers, homeCode, awayCode, row.Url)
	// processDocAndTeam(res, &awayPlayers, awayCode)

	processedResult := DetailedMatchContainer{
		matchUrl: row.Url,
		HomeTeamOnlineID: homeCode,
		AwayTeamOnlineID: awayCode,
		homePlayers: homePlayers,
		awayPlayers: awayPlayers,
		isknockout: true,
	}

	fmt.Printf("== Successfully finished scraping url: %s\n", row.Url)
	c <- processedResult
}

func parseLeagueResponse(res *http.Response, c chan DetailedMatchContainer, row database.GetLeagueMatchUrlsAndTeamOnlineIdsWOffsetRow) {

	homeCode := strings.Split(row.HomeTeamOnlineID,"/")[3]
	awayCode := strings.Split(row.AwayTeamOnlineID,"/")[3]

	var homePlayers []PlayerDetailContainer
	var awayPlayers []PlayerDetailContainer

	defer res.Body.Close()
	processDocAndTeam(res, &homePlayers, &awayPlayers, homeCode, awayCode, row.Url)

	processedResult := DetailedMatchContainer{
		matchUrl: row.Url,
		HomeTeamOnlineID: homeCode,
		AwayTeamOnlineID: awayCode,
		homePlayers: homePlayers,
		awayPlayers: awayPlayers,
		isknockout: false,
	}

	fmt.Printf("== Successfully finished scraping url: %s\n", row.Url)
	c <- processedResult
}
