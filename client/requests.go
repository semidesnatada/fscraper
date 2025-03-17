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

const (
	baseUrl = "https://fbref.com/en/comps/"
	oldPremUrl = "9/2015-2016/schedule/2015-2016-Premier-League-Scores-and-Fixtures"
	premUrl = "9/2023-2024/schedule/2023-2024-Premier-League-Scores-and-Fixtures"
	ligaUrl = "12/2022-2023/schedule/2022-2023-La-Liga-Scores-and-Fixtures"
	vOldPremUrl = "9/1991-1992/schedule/1991-1992-Premier-League-Scores-and-Fixtures"
)

type compMetaRecord struct {
	Name, OnlineCode string
	EarliestYear, LatestYear int
}

func GenerateCompsForSearching() []CompetitionSeasonSummary {
	
	// codes["Champions-League"] = "8"
	// codes["FA-Cup"] = "514"

	var output []CompetitionSeasonSummary
	LeagueParameters := getLeagueParams()

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
	matches := ParseLeagueResults(*res, *comp)
	comp.Data = matches
	// PrintMatches(matches, 5)
	return nil
}

func ScrapeLeagues(s *config.State) error {
	comps := GenerateCompsForSearching()

	// comps := [1]CompetitionSeasonSummary{}
	// comps[0] = CompetitionSeasonSummary{
	// 	CompetitionName: "Ligue-1",
	// 	CompetitionSeason: "2019-2020",
	// 	CompetitionOnlineID: "13",
	// 	Url:"https://fbref.com/en/comps/13/2019-2020/schedule/2019-2020-Ligue-1-Scores-and-Fixtures",
	// }

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
				fmt.Println("===================================================")
				fmt.Printf("%s, season: %s already exists in database. moving on to next scrape.\n", comp.CompetitionName, comp.CompetitionSeason)
				fmt.Println("===================================================")
				continue
			}
			//if doesn't exist in db, then scrape it.
			go GoScrape(comp, channel)
			<- ticker.C
		}
		close(channel)
	}(resultsChannel)
	
	for matches := range resultsChannel{
		fmt.Println("===================================================")
		fmt.Printf("Started storing: %s, season: %s\n", matches.CompetitionName, matches.CompetitionSeason)
		storeErr := StoreMatchSummaries(s, matches)
		if storeErr != nil {
			fmt.Println(storeErr.Error())
			os.Exit(1)
		}
		fmt.Printf("Successfully concluded storing: %s, season: %s\n", matches.CompetitionName, matches.CompetitionSeason)
		fmt.Println("===================================================")
	}
	return nil
}

func GoScrape(comp CompetitionSeasonSummary, channel chan CompetitionSeasonSummary) {
	
	fmt.Println("===================================================")
	fmt.Printf("Started scraping: %s, %s\n", comp.CompetitionName, comp.CompetitionSeason)
	err := ScrapeLeagueFromUrl(&comp)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("Successfully finished scraping: %s, %s\n", comp.CompetitionName, comp.CompetitionSeason)
	fmt.Println("===================================================")
	channel <- comp
}

func getLeagueParams() []compMetaRecord {
	return []compMetaRecord {
	{
		Name: "Premier-League",
		OnlineCode: "9",
		EarliestYear: 1991,
		LatestYear: 2025,
	},
	{
		Name: "La-Liga",
		OnlineCode: "12",
		EarliestYear: 1991,
		LatestYear: 2025,
	},
	{
		Name: "Bundesliga",
		OnlineCode: "20",
		EarliestYear: 1991,
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
		EarliestYear: 1991,
		LatestYear: 2025,
	},
	{
		Name: "Championship",
		OnlineCode: "10",
		EarliestYear: 2014,
		LatestYear: 2025,
}}}