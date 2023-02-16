package types

import (
	"database/sql/driver"
	"encoding/json"
)

type RiotMatchRes struct {
	MetaData struct {
		DataVersion  string   `validate:"required" json:"data_version"`
		MatchId      string   `validate:"required" json:"match_id"    `
		Participants []string `validate:"required" json:"participants"`
	} `validate:"required" json:"metadata"`
	Info MatchInfo `validate:"required" json:"info"`
}

type MatchInfo struct {
	Date      int64   `validate:"required" json:"game_datetime"`
	Length    float32 `validate:"required" json:"game_length"`
	Version   string  `validate:"required" json:"game_version"`
	QueueId   int     `validate:"required" json:"queue_id"`
	GameType  string  `validate:"required" json:"tft_game_type"`
	SetName   string  `validate:"required" json:"tft_set_core_name"`
	SetNumber int     `validate:"required" json:"tft_set_number"`
	Comps     []struct {
		Augments  []string `validate:"required" json:"augments"`
		Companion struct {
			ContentId string `validate:"required" json:"content_ID"`
			ItemId    int    `validate:"required" json:"item_ID"`
			SkinId    int    `validate:"required" json:"skin_ID"`
			Species   string `validate:"required" json:"species"`
		} `validate:"required" json:"companion"`
		RemainingGold     int     `validate:"required" json:"gold_left"`
		RoundEliminated   int     `validate:"required" json:"last_round"`
		Level             int     `validate:"required" json:"level"`
		Placement         int     `validate:"required" json:"placement"`
		PlayersEliminated int     `validate:"required" json:"players_eliminated"`
		Puuid             string  `validate:"required" json:"puuid"`
		TimeEliminated    float32 `validate:"required" json:"time_eliminated"`
		DamageDealt       int     `validate:"required" json:"total_damage_to_players"`
		Traits            []struct {
			Name       string `validate:"required" json:"name"`
			NumUnits   int    `validate:"required" json:"num_units"`
			Style      int    `validate:"required" json:"style"`
			TierActive int    `validate:"required" json:"tier_current"`
			TierMax    int    `validate:"required" json:"tier_total"`
		} `validate:"required" json:"traits"`
		Units []struct {
			Id        string   `validate:"required" json:"character_id"`
			ItemNames []string `validate:"required" json:"itemNames"`
			ItemIds   []int    `validate:"required" json:"items"`
			Name      string   `validate:"required" json:"name"`
			Rarity    int      `validate:"required" json:"rarity"`
			Tier      int      `validate:"required" json:"tier"`
		} `validate:"required" json:"units"`
	} `validate:"required" json:"participants"`
}

func (info MatchInfo) Value() (driver.Value, error) {
	return json.Marshal(info)
}
