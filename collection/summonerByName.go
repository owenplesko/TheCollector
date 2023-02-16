package collection

import (
	"fmt"
	db "thecollector/database"
	"thecollector/riot"
	"thecollector/util"
)

type SummonerByNameCollecter struct {
	RawName string
	Region  string
}

func NewSummonerByNameCollecter(region string, name string) SummonerByNameCollecter {
	return SummonerByNameCollecter{
		RawName: util.ToRawName(name),
		Region:  region,
	}
}

func (c SummonerByNameCollecter) Id() string {
	return c.Region + c.RawName
}

func (c SummonerByNameCollecter) Collect() error {
	fmt.Printf("Collecting summoner %v\n", c.RawName)
	// get summoner from riot
	summoner, err := riot.GetSummonerByName(c.Region, c.RawName)
	if err != nil {
		fmt.Printf("Error getting summoner %s from riot: %s\n", c.RawName, err)
		return err
	}

	// store summoner in db
	err = db.UpsertSummoner(summoner, c.Region)
	if err != nil {
		fmt.Printf("Error inserting summoner %s into db %s\n", c.RawName, err)
	}
	return err
}
