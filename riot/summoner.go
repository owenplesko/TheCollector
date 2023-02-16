package riot

import (
	"fmt"
	"thecollector/types"
)

func GetSummonerByName(region string, name string) (*types.RiotSummonerRes, error) {
	summoner := new(types.RiotSummonerRes)
	url := fmt.Sprintf("https://%v.api.riotgames.com/tft/summoner/v1/summoners/by-name/%v?api_key=%v", region, name, key)
	err := getJson(url, summoner)
	if err != nil {
		return summoner, err
	}
	err = validate.Struct(summoner)
	return summoner, err
}

func GetSummonerByPuuid(region string, puuid string) (*types.RiotSummonerRes, error) {
	summoner := new(types.RiotSummonerRes)
	url := fmt.Sprintf("https://%v.api.riotgames.com/tft/summoner/v1/summoners/by-puuid/%v?api_key=%v", region, puuid, key)
	err := getJson(url, summoner)
	if err != nil {
		return summoner, err
	}
	err = validate.Struct(summoner)
	return summoner, err
}
