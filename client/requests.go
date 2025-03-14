package client

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/semidesnatada/fscraper/config"
)

const (
	baseUrl = "https://fbref.com/en/comps/"
	oldPremUrl = "9/2015-2016/schedule/2015-2016-Premier-League-Scores-and-Fixtures"
	premUrl = "9/2023-2024/schedule/2023-2024-Premier-League-Scores-and-Fixtures"
	ligaUrl = "12/2022-2023/schedule/2022-2023-La-Liga-Scores-and-Fixtures"
	vOldPremUrl = "9/1991-1992/schedule/1991-1992-Premier-League-Scores-and-Fixtures"
)

func GenerateCompsForSearching() []CompetitionSeasonSummary {

	leagueNames := []string{}
	leagueNames = append(leagueNames, "Premier-League")
	leagueNames = append(leagueNames, "La-Liga")

	codes := make(map[string]string)
	codes["Premier-League"] = "9"
	codes["La-Liga"] = "12"

	// years := []string{"1991-1992","1992-1993","1993-1994","1994-1995","1995-1996","1996-1997","1997-1998",
	// 				"1998-1999","1999-2000","2000-2001","2001-2002","2002-2003","2003-2004","2004-2005",
	// 				"2005-2006","2006-2007","2007-2008","2008-2009","2009-2010","2010-2011","2011-2012",
	// 				"2012-2013","2013-2014","2014-2015","2015-2016","2016-2017","2017-2018","2018-2019",
	// 				"2019-2020","2020-2021","2021-2022","2022-2023","2023-2024"}

	years := []string{"2014-2015","2015-2016","2016-2017","2017-2018","2018-2019",
					"2019-2020","2020-2021","2021-2022","2022-2023","2023-2024"}
	
	var output []CompetitionSeasonSummary

	for _, year := range years {
		for _, league := range leagueNames {
			output = append(output,
				CompetitionSeasonSummary{
					CompetitionName: league,
					CompetitionSeason: year,
					Url: fmt.Sprintf("%s%s/%s/schedule/%s-%s-Scores-and-Fixtures", baseUrl, codes[league], year, year, league),
				},
			)
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
	// PrintMatches(parsedResult.Data, 5)
	return nil
}

func ScrapeLeagues(s *config.State) {
	comps := GenerateCompsForSearching()

	// urls := []string{"https://fbref.com/en/comps/12/2015-2016/schedule/2015-2016-La-Liga-Scores-and-Fixtures"}

	resultsChannel := make(chan CompetitionSeasonSummary)

	go func(channel chan CompetitionSeasonSummary) {
		ticker := time.NewTicker(time.Second * 6)
		for _, comp := range comps {
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