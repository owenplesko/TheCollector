package types

type Summoner struct {
	Puuid              string `json:"puuid"`
	Region             string `json:"region"`
	SummonerId         string `json:"summonerId"`
	AccountId          string `json:"accountId"`
	DisplayName        string `json:"displayName"`
	RawName            string `json:"rawName"`
	ProfileIconId      int    `json:"profileIconId"`
	RevisionDate       int    `json:"revisionDate"`
	SummonerLevel      int    `json:"summonerLevel"`
	LastUpdated        int64  `json:"lastUpdated"`
	RankLastUpdated    int64  `json:"rankLastUpdated"`
	MatchesLastUpdated int64  `json:"matchesLastUpdated"`
}
