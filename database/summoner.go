package db

import (
	"strings"
	"thecollector/riot"
	"time"
)

func UpsertSummoner(summoner *riot.Summoner) error {
	query := `
		INSERT INTO summoner (
			puuid,
			region,
			summoner_id,
			account_id,
			profile_icon_id,
			revision_date,
			display_name,
			raw_name,
			summoner_level,
			last_updated
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (puuid) DO UPDATE SET
			region= $2,
			summoner_id= $3,
			account_id= $4,
			profile_icon_id= $5,
			revision_date= $6,
			display_name= $7,
			raw_name= $8,
			summoner_level= $9,
			last_updated= $10;
	`

	_, err := db.Exec(query,
		summoner.Puuid,
		summoner.Region,
		summoner.Id,
		summoner.AccountId,
		summoner.ProfileIconId,
		summoner.RevisionDate,
		summoner.Name,
		toRawName(summoner.Name),
		summoner.SummonerLevel,
		time.Now().Unix(),
	)
	return err
}

func SummonerExists(puuid string) bool {
	query := `SELECT EXISTS (SELECT 1 FROM summoner WHERE puuid=$1);`
	var exists bool
	row := db.QueryRow(query, puuid)
	row.Scan(&exists)
	return exists
}

func UpdateMatchesLastUpdated(puuid string, time int64) error {
	query := `UPDATE summoner SET matches_last_updated=$1 WHERE puuid=$2;`
	_, err := db.Exec(query, time, puuid)
	return err
}

func QueryPuuidMatchesLastUpdated(exclude []string) (string, int64, error) {
	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT puuid, matches_last_updated FROM summoner ")
	if len(exclude) > 0 {
		queryBuilder.WriteString("WHERE puuid NOT IN ('")
		queryBuilder.WriteString(strings.Join(exclude, "', '"))
		queryBuilder.WriteString("') ")
	}
	queryBuilder.WriteString("ORDER BY matches_last_updated LIMIT 1")
	var puuid string
	var lastUpdated int64
	row := db.QueryRow(queryBuilder.String())
	err := row.Scan(&puuid, &lastUpdated)
	return puuid, lastUpdated, err
}

func toRawName(displayName string) string {
	return strings.ToLower(strings.TrimSpace(displayName))
}
