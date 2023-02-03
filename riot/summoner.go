package riot

import "fmt"

type Summoner struct {
	Puuid         string `json:"puuid"         validate:"required"`
	Region        string `                     validate:"required"`
	Id            string `json:"id"            validate:"required"`
	AccountId     string `json:"accountId"     validate:"required"`
	Name          string `json:"name"          validate:"required"`
	ProfileIconId int    `json:"profileIconId"`
	RevisionDate  int    `json:"revisionDate"  validate:"required"`
	SummonerLevel int    `json:"summonerLevel" validate:"required"`
}

func GetSummonerByName(region string, name string) (*Summoner, error) {
	summoner := new(Summoner)
	summoner.Region = region
	url := fmt.Sprintf("https://%v.api.riotgames.com/tft/summoner/v1/summoners/by-name/%v?api_key=%v", region, name, key)
	err := getJson(url, summoner)
	if err != nil {
		return summoner, err
	}
	err = validate.Struct(summoner)
	return summoner, err
}

func GetSummonerByPuuid(region string, puuid string) (*Summoner, error) {
	summoner := new(Summoner)
	summoner.Region = region
	url := fmt.Sprintf("https://%v.api.riotgames.com/tft/summoner/v1/summoners/by-puuid/%v?api_key=%v", region, puuid, key)
	err := getJson(url, summoner)
	if err != nil {
		return summoner, err
	}
	err = validate.Struct(summoner)
	return summoner, err
}
