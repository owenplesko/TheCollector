package collect

import (
	"fmt"
	db "thecollector/database"
	"thecollector/riot"
)

type SummonerByNameCollecter struct {
	RawName string
	Region  string
}

func NewSummonerByNameCollecter(region string, rawName string) SummonerByNameCollecter {
	return SummonerByNameCollecter{
		rawName,
		region,
	}
}

func (c SummonerByNameCollecter) Id() string {
	return c.Region + c.RawName
}

func (c SummonerByNameCollecter) Collect(_ bool) error {
	fmt.Printf("Collecting summoner %v\n", c.RawName)
	// get summoner from riot
	summoner, err := riot.GetSummonerByName(c.Region, c.RawName)
	if err != nil {
		fmt.Printf("Error getting summoner %s from riot: %s\n", c.RawName, err)
		return err
	}

	// store summoner in db
	err = db.UpsertSummoner(summoner)
	if err != nil {
		fmt.Printf("Error inserting summoner %s into db %s\n", c.RawName, err)
	}
	return err
}
