package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/semidesnatada/fscraper/client"
	"github.com/semidesnatada/fscraper/config"

	_ "github.com/lib/pq"
)


func main() {
	
	const DB_URL = "postgres://seanlowery:@localhost:5432/fscraped?sslmode=disable"
	s := config.CreateState(DB_URL)
	s.DeleteDatabases()
	
	urls := client.GenerateUrlsForSearching()

	ticker := time.NewTicker(time.Second * 5)

	for _, url := range urls {

		matches, err := testScrapingLeague(url)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		storeErr := client.StoreMatchSummaries(&s, matches)
		
		if storeErr != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		<- ticker.C

	}
	fmt.Println("finished testing database")

	os.Exit(0)
}


func testScrapingLeague(url string) (client.CompetitionSeasonSummary, error) {

	res, err := http.Get(url)
	
	if err != nil {
		return client.CompetitionSeasonSummary{}, err
	}

	defer res.Body.Close()
	parsedResult := client.ParseLeagueResults(*res)
	client.PrintMatches(parsedResult.Data, 5)

	return parsedResult, nil
}

// func main () {

	// testDB()


	// fmt.Println("Hello worldino")
// }


// func testDB () {

// 	const DB_URL = "postgres://seanlowery:@localhost:5432/fscraped?sslmode=disable"

// 	db, err := sql.Open("postgres", DB_URL)
// 	if err != nil {
// 		fmt.Printf("Programme terminated: %s\n", err.Error())
// 		os.Exit(1)
// 	}
// 	dbQueries := database.New(db)


// 	_, insertErr := dbQueries.CreateMatch(context.Background(),
// 							database.CreateMatchParams{
// 								ID: uuid.New(),
// 								HomeTeam: "Jotetnham Flotspurg",
// 								AwayTeam: "The Players",
// 								HomeGoals: 2,
// 								AwayGoals: 7,
// 							})
// 	if insertErr != nil {
// 		fmt.Printf("Programme terminated: %s\n", insertErr.Error())
// 		os.Exit(1)
// 	}

// 	matches, queryErr := dbQueries.GetMatches(context.Background())
// 	if queryErr != nil {
// 		fmt.Printf("Programme terminated: %s\n", queryErr.Error())
// 		os.Exit(1)
// 	}

// 	for _, match := range matches {
// 		fmt.Println("=================")
// 		fmt.Printf("Match of the Day: %s vs %s\n", match.HomeTeam, match.AwayTeam)
// 		fmt.Printf("Final Result: %d : %d\n", match.HomeGoals, match.AwayGoals)
// 	}

// }