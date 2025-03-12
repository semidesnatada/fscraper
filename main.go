package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/semidesnatada/fscraper/database"

	_ "github.com/lib/pq"
)

func main () {

	const DB_URL = "postgres://seanlowery:@localhost:5432/fscraped?sslmode=disable"

	db, err := sql.Open("postgres", DB_URL)
	if err != nil {
		fmt.Printf("Programme terminated: %s\n", err.Error())
		os.Exit(1)
	}
	dbQueries := database.New(db)


	_, insertErr := dbQueries.CreateMatch(context.Background(),
							database.CreateMatchParams{
								ID: uuid.New(),
								HomeTeam: "Jotetnham Flotspurg",
								AwayTeam: "The Players",
								HomeGoals: 2,
								AwayGoals: 7,
							})
	if insertErr != nil {
		fmt.Printf("Programme terminated: %s\n", insertErr.Error())
		os.Exit(1)
	}

	matches, queryErr := dbQueries.GetMatches(context.Background())
	if queryErr != nil {
		fmt.Printf("Programme terminated: %s\n", queryErr.Error())
		os.Exit(1)
	}

	for _, match := range matches {
		fmt.Println("=================")
		fmt.Printf("Match of the Day: %s vs %s\n", match.HomeTeam, match.AwayTeam)
		fmt.Printf("Final Result: %d : %d\n", match.HomeGoals, match.AwayGoals)
	}

	fmt.Println("Hello worldino")
}
