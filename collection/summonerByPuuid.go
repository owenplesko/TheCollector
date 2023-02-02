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
		return err
	}

	// store summoner in db
	return db.UpsertSummoner(summoner)
}
