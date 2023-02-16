package db

import (
	"encoding/json"
	"fmt"
	"strings"
	"thecollector/types"
)

func StoreMatch(match *types.RiotMatchRes) error {
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
	db.QueryRow(`SELECT EXISTS (SELECT 1 FROM match WHERE id=$1)`, matchId).Scan(&exists)
	return exists
}

func QueryMatchHistory(puuid string, count int) ([]*types.Match, error) {
	var matchHistory []*types.Match

	rows, err := db.Query(`
		SELECT match.id, match.date, match.data_version, match.data
		FROM match JOIN summoner_match ON match.id = summoner_match.match_id
		WHERE summoner_match.summoner_puuid = $1
		ORDER BY match.date DESC
		LIMIT $2
	`, puuid, count)

	if err != nil {
		return matchHistory, err
	}
	defer rows.Close()

	// scan each row into a Match struct and add to the slice of matches
	for rows.Next() {
		match := new(types.Match)
		var jsonData []byte
		if err := rows.Scan(&match.Id, &match.Date, &match.DataVersion, &jsonData); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(jsonData, &match.Info); err != nil {
			return nil, err
		}
		matchHistory = append(matchHistory, match)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matchHistory, nil
}
