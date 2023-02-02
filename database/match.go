package db

import (
	"encoding/json"
	"fmt"
	"strings"
	"thecollector/riot"
)

func StoreMatch(match *riot.Match) error {
	// start transaction for new match and summoner_match connections
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// on err rollback transaction to prevent dangling matches
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// insert match
	query := `
		INSERT INTO match (
			id,
			date,
			data_version,
			data
		)
		VALUES (?, ?, ?, ?)
	`
	data, err := json.Marshal(match.Info)
	if err != nil {
		return err
	}
	_, err = tx.Exec(query,
		match.MetaData.MatchId,
		match.Info.Date,
		match.MetaData.DataVersion,
		data,
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
	if err != nil {
		return err
	}

	return tx.Commit()
}

func MatchExists(matchId string) bool {
	query := `SELECT EXISTS (SELECT 1 FROM match WHERE id=?);`
	var exists bool
	row := db.QueryRow(query, matchId)
	row.Scan(&exists)
	return exists
}
