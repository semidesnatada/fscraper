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

func GenerateKnockoutsForSearching() []CompetitionSeasonSummary {

	var output []CompetitionSeasonSummary
	KnockoutParameters := GetKnockoutParams()

	for _, comp := range KnockoutParameters {
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

func ScrapeKnockoutFromUrl(comp *CompetitionSeasonSummary) (error) {

	res, err := http.Get(comp.Url)
	
	if err != nil {
		return err
	}

	defer res.Body.Close()
	matches := parseKnockoutResults(*res, *comp)
	comp.Data = matches
	// PrintKnockoutMatches(matches, 5)
	return nil
}

func ScrapeKnockouts(s *config.State) error {
	//get urls to scrape
	comps := GenerateKnockoutsForSearching()

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
				fmt.Println("error in checking whether competition exists in db in client.requests")
				os.Exit(1)
			}

			if alreadyExists {
				fmt.Println("=====")
				fmt.Printf("%s, season: %s already exists in database. moving on to next scrape.\n", comp.CompetitionName, comp.CompetitionSeason)
				fmt.Println("=====")
				continue
			}
			//if doesn't exist in db, then scrape it.
			GoScrapeKnockout(comp, channel)
			<- ticker.C
		}
		close(channel)
	}(resultsChannel)
	
	// once data is scraped, store it in the database.
	for matches := range resultsChannel{
		fmt.Printf("=== Started storing: %s, season: %s\n", matches.CompetitionName, matches.CompetitionSeason)
		storeErr := storeKnockoutMatchSummaries(s, matches)
		if storeErr != nil {
			fmt.Println(storeErr.Error())
			os.Exit(1)
		}
		fmt.Printf("==== Successfully concluded storing: %s, season: %s\n", matches.CompetitionName, matches.CompetitionSeason)
		fmt.Println("__________")
	}
	return nil
}

func GoScrapeKnockout(comp CompetitionSeasonSummary, channel chan CompetitionSeasonSummary) {
	
	fmt.Printf("= Started scraping: %s, %s\n", comp.CompetitionName, comp.CompetitionSeason)
	err := ScrapeKnockoutFromUrl(&comp)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("== Successfully finished scraping: %s, %s\n", comp.CompetitionName, comp.CompetitionSeason)
	
	channel <- comp
}

func GetKnockoutParams() []compKnockoutMetaRecord {
	return []compKnockoutMetaRecord {
	{
		Name: "FA-Cup",
		OnlineCode: "514",
		EarliestYear: 1992,
		LatestYear: 2024,
	},
	{
		Name: "Champions-League",
		OnlineCode: "8",
		EarliestYear: 1992,
		LatestYear: 2024,
	},
	{
		Name: "Europa-League",
		OnlineCode: "19",
		EarliestYear: 1992,
		LatestYear: 2024,
	},
}}