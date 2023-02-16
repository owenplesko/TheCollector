package types

type RiotSummonerRes struct {
	Puuid         string `json:"puuid"       validate:"required"`
	SummonerId    string `json:"id"          validate:"required"`
	AccountId     string `json:"accountId"   validate:"required"`
	Name          string `json:"name"        validate:"required"`
	ProfileIconId int    `json:"profileIconId"`
	RevisionDate  int    `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}
