package client

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/semidesnatada/fscraper/config"
	"github.com/semidesnatada/fscraper/database"
)

// brings together the scraping functions from "Parse" file, and processing and DB storage functions from "Store" file.

func GenerateLeaguesForSearching() []CompetitionSeasonSummary {

	var output []CompetitionSeasonSummary
	LeagueParameters := GetLeagueParams()

	for _, comp := range LeagueParameters {
		for i := comp.EarliestYear; i < comp.LatestYear; i++ {
			season := fmt.Sprintf("%d-%d",i,i+1)
			x := CompetitionSeasonSummary{
				CompetitionName: comp.Name,
				CompetitionSeason: season,
				CompetitionOnlineID: comp.OnlineCode,
				Url: fmt.Sprintf("%s%s/%s/schedule/%s-%s-Scores-and-Fixtures", baseUrl, comp.OnlineCode, season, season, comp.Name),
			}
			output = append(output, x)
		}
	}

	return output
}

func ScrapeLeagueFromUrl(comp *CompetitionSeasonSummary) (error) {

	res, err := http.Get(comp.Url)
	
	if err != nil {
		return err
	}

	defer res.Body.Close()
	matches := parseLeagueResults(*res, *comp)
	comp.Data = matches
	// PrintLeagueMatches(matches, 5)
	return nil
}

func ScrapeLeagues(s *config.State) error {
	//get urls to scrape
	comps := GenerateLeaguesForSearching()

	//create channel on which to return structured data
	resultsChannel := make(chan CompetitionSeasonSummary)

	go func(channel chan CompetitionSeasonSummary) {
		ticker := time.NewTicker(time.Second * 6)
		for _, comp := range comps {
			//check if comp already exists in db

			alreadyExists, checkErr := s.DB.CheckIfCompetitionExistsByNameAndSeason(
				context.Background(),
				database.CheckIfCompetitionExistsByNameAndSeasonParams{
					Name: comp.CompetitionName,
					Season: comp.CompetitionSeason,
				},
			)
			if checkErr != nil {
				fmt.Println("error in checking whether league exists in db in client.requests")
				os.Exit(1)
			}

			if alreadyExists {
				fmt.Println("=====")
				fmt.Printf("%s, season: %s already exists in database. moving on to next scrape.\n", comp.CompetitionName, comp.CompetitionSeason)
				fmt.Println("=====")
				continue
			}
			//if doesn't exist in db, then scrape it.
			GoScrapeLeague(comp, channel)
			<- ticker.C
		}
		close(channel)
	}(resultsChannel)
	
	// once data is scraped, store it in the database.
	for matches := range resultsChannel{
		fmt.Printf("=== Started storing: %s, season: %s\n", matches.CompetitionName, matches.CompetitionSeason)
		storeErr := storeLeagueMatchSummaries(s, matches)
		if storeErr != nil {
			fmt.Println(storeErr.Error())
			os.Exit(1)
		}
		fmt.Printf("==== Successfully concluded storing: %s, season: %s\n", matches.CompetitionName, matches.CompetitionSeason)
		fmt.Println("__________")
	}
	return nil
}

func GoScrapeLeague(comp CompetitionSeasonSummary, channel chan CompetitionSeasonSummary) {
	
	fmt.Printf("= Started scraping: %s, %s\n", comp.CompetitionName, comp.CompetitionSeason)
	err := ScrapeLeagueFromUrl(&comp)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("== Successfully finished scraping: %s, %s\n", comp.CompetitionName, comp.CompetitionSeason)
	
	channel <- comp
}

func GetLeagueParams() []compLeagueMetaRecord {
	return []compLeagueMetaRecord {
	{
		Name: "Premier-League",
		OnlineCode: "9",
		EarliestYear: 1992,
		LatestYear: 2025,
	},
	{
		Name: "La-Liga",
		OnlineCode: "12",
		EarliestYear: 1992,
		LatestYear: 2025,
	},
	{
		Name: "Bundesliga",
		OnlineCode: "20",
		EarliestYear: 1992,
		LatestYear: 2025,
	},
	{
		Name: "Ligue-1",
		OnlineCode: "13",
		EarliestYear: 1995,
		LatestYear: 2025,
	},
	{
		Name: "Serie-A",
		OnlineCode: "11",
		EarliestYear: 1992,
		LatestYear: 2025,
	},
// 	{
// 		Name: "Championship",
// 		OnlineCode: "10",
// 		EarliestYear: 2014,
// 		LatestYear: 2025,
// },
}}