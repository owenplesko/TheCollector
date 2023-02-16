package riot

import (
	"fmt"
	"thecollector/types"
)

type MatchHistory struct {
	Matches []string
}

func GetMatchDetails(matchId string) (*types.RiotMatchRes, error) {
	match := new(types.RiotMatchRes)
	region := "americas" // TODO: select based on region prefix
	url := fmt.Sprintf("https://%v.api.riotgames.com/tft/match/v1/matches/%v?api_key=%v", region, matchId, key)
	err := getJson(url, match)
	if err != nil {
		return match, err
	}
	err = validate.Struct(match)
	return match, err
}

func GetMatchHistory(puuid string, after int64) ([]string, error) {
	var history []string
	region := "americas"
	count := 200
	if matchesAfter > after {
		after = matchesAfter
	}

	url := fmt.Sprintf("https://%v.api.riotgames.com/tft/match/v1/matches/by-puuid/%v/ids?startTime=%v&count=%v&api_key=%v", region, puuid, after, count, key)
	err := getJson(url, &history)
	return history, err
}
