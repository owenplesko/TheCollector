package collect

import (
	"fmt"
	db "thecollector/database"
	"thecollector/riot"
)

type SummonerByPuuidCollecter struct {
	Puuid  string
	Region string
}

func NewSummonerByPuuidCollecter(region string, puuid string) SummonerByPuuidCollecter {
	return SummonerByPuuidCollecter{
		puuid,
		region,
	}
}

func (c SummonerByPuuidCollecter) Id() string {
	return c.Puuid
}

func (c SummonerByPuuidCollecter) Collect() error {
	fmt.Printf("Collecting summoner %v\n", c.Puuid)
	// get summoner from riot
	summoner, err := riot.GetSummonerByPuuid(c.Region, c.Puuid)
	if err != nil {
		fmt.Printf("Error getting summoner %s from riot: %s\n", c.Puuid, err)
		return err
	}

	// store summoner in db
	err = db.UpsertSummoner(summoner)
	if err != nil {
		fmt.Printf("Error inserting summoner %s into db %s\n", c.Puuid, err)
	}
	return err
}
