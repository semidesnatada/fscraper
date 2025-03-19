package client

const baseUrl = "https://fbref.com/en/comps/"
const baseMatchUrl = "https://fbref.com"

type MatchSummary struct {
	knockout bool
	data map[string]string
}

type CompetitionSeasonSummary struct {
	Data []MatchSummary
	CompetitionName string
	CompetitionSeason string
	CompetitionOnlineID string
	Knockout bool
	Url string
}

type compLeagueMetaRecord struct {
	Name, OnlineCode string
	EarliestYear, LatestYear int
}
type compKnockoutMetaRecord struct {
	Name, OnlineCode, TablesToScrape string
	EarliestYear, LatestYear int
}