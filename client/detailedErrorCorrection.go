package client

import (
	"context"
	"fmt"

	"github.com/semidesnatada/fscraper/config"
	"github.com/semidesnatada/fscraper/database"
)


func DeleteAndRescrapeLeagueMatches(s *config.State) error {

	urls, err := s.DB.GetUrlsToRescrape(context.Background())
	if err != nil {
		return err
	}

	// for _, url := range urls {
	// 	fmt.Printf("url is: %s\n",url.Url)
	// }

	fmt.Println()
	fmt.Printf("number of urls to rescrape: %d\n", len(urls))
	fmt.Println()
	fUrls := ConvertToCorrectFormat(urls)

	// for testing
	// fUrls := []database.GetLeagueMatchUrlsAndTeamOnlineIdsWOffsetRow{}
	// next := database.GetLeagueMatchUrlsAndTeamOnlineIdsWOffsetRow{}

	//test2
	// next.Url = "https://fbref.com/en/matches/829fd962/Athletic-Club-Real-Betis-March-13-2016-La-Liga"
	// next.HomeTeamOnlineID = "/en/squads/2b390eca/2015-2016/Athletic-Club-Stats"
	// next.AwayTeamOnlineID = "https://fbref.com/en/squads/fc536746/2015-2016/Real-Betis-Stats"

	//test1
	// next.Url = "/en/matches/474ce5b7/Schalke-04-Bayer-Leverkusen-April-9-1996-Bundesliga"
	// next.HomeTeamOnlineID = "/en/squads/c539e393/1995-1996/Schalke-04-Stats"
	// next.AwayTeamOnlineID = "/en/squads/c7a9f859/1995-1996/Bayer-Leverkusen-Stats"
	// fUrls = append(fUrls, next)
	
	dErr := DeleteAndReScrapeLeagueMatchByUrl(s, fUrls)
	if dErr != nil {
		return dErr
	}

	return nil
}


func DeleteAndReScrapeLeagueMatchByUrl(s *config.State, matchUrl []database.GetLeagueMatchUrlsAndTeamOnlineIdsWOffsetRow) error {

	fmt.Println("Started deleting")
	for i, item := range matchUrl {
		err := s.DB.DeleteRecordsForGivenMatchUrl(context.Background(), item.Url)
		if err != nil {
			return fmt.Errorf("error deleting the record %s, details: %s", item.Url, err.Error())
		}
		fmt.Printf("Deleted record %d, url: %s\n",i,item.Url)
	}
	fmt.Println("Completed deleting")
	
	processedMatchChannel := make(chan DetailedMatchContainer)

	fmt.Println("Started processing all urls")
	go processLeagueUrlBlock(s, matchUrl, processedMatchChannel)

	for processed := range processedMatchChannel {
		fmt.Printf("=== Started storing url: %s \n", processed.matchUrl)
		storeDetailedMatches(s, processed)
	}
	return nil
}

func ConvertToCorrectFormat(urls_in []database.GetUrlsToRescrapeRow) []database.GetLeagueMatchUrlsAndTeamOnlineIdsWOffsetRow {

	urls := []database.GetLeagueMatchUrlsAndTeamOnlineIdsWOffsetRow{}

	for _, item := range urls_in {
		next := database.GetLeagueMatchUrlsAndTeamOnlineIdsWOffsetRow{}
		next.Url = item.Url
		next.HomeTeamOnlineID = item.HomeTeamOnlineID
		next.AwayTeamOnlineID = item.AwayTeamOnlineID
		urls = append(urls, next)
	}
	return urls
}