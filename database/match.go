package db

import (
	"fmt"
	"strings"
	"thecollector/riot"
)

func StoreMatch(match *riot.Match) error {
	// insert match
	query := `INSERT INTO match (id, date, data_version, data) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(query,
		match.MetaData.MatchId,
		match.Info.Date,
		match.MetaData.DataVersion,
		match.Info,
	)
	if err != nil {
		return err
	}

	// insert summoner_match connections
	var queryBuilder strings.Builder
	queryBuilder.WriteString("INSERT INTO summoner_match (summoner_puuid, match_id) VALUES ")
	for i, puuid := range match.MetaData.Participants {
		if i != 0 {
			queryBuilder.WriteString(",")
		}
		queryBuilder.WriteString(fmt.Sprintf("('%s', '%s')", puuid, match.MetaData.MatchId))
	}
	_, err = db.Exec(queryBuilder.String())

	return err
}

func MatchExists(matchId string) bool {
	var exists bool
	err := db.QueryRow(`SELECT EXISTS (SELECT 1 FROM match WHERE id=$1)`, matchId).Scan(&exists)
	if err != nil {
		fmt.Printf("Error checking if match %s exists: %s\n", matchId, err)
		return false
	}
	return exists
}
