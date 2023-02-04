package db

import (
	"fmt"
	"strings"
	"thecollector/riot"
)

func StoreMatch(match *riot.Match) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// insert match
	_, err = tx.Exec(`INSERT INTO match (id, date, data_version, data) VALUES ($1, $2, $3, $4)`,
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
	queryBuilder.WriteString(`INSERT INTO summoner_match (summoner_puuid, match_id) VALUES `)
	for i, puuid := range match.MetaData.Participants {
		if i != 0 {
			queryBuilder.WriteString(",")
		}
		queryBuilder.WriteString(fmt.Sprintf("('%s', '%s')", puuid, match.MetaData.MatchId))
	}
	_, err = tx.Exec(queryBuilder.String())

	return tx.Commit()
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
