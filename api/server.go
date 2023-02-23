package api

import (
	"errors"
	"net/url"
	"thecollector/collection"
	db "thecollector/database"

	"github.com/gofiber/fiber/v2"
)

// summoner/:puuid
func SummonerByPuuid(c *fiber.Ctx) error {
	puuid := c.Params("puuid")
	if !db.SummonerPuuidExists(puuid) {
		c.SendStatus(404)
		return errors.New(puuid + " not in db")
	}
	summoner, err := db.QuerySummonerByPuuid(puuid)
	if err != nil {
		c.Status(404).SendString(err.Error())
		return err
	}
	c.Status(200).JSON(*summoner)
	return nil
}

// summoner/:region/:name
func SummonerByName(c *fiber.Ctx) error {
	region := c.Params("region")
	name, err := url.QueryUnescape(c.Params("name"))
	if err != nil {
		c.SendStatus(400)
	}
	if !db.SummonerNameExists(region, name) {
		<-prioritySummonerScheduler.Schedule(collection.NewSummonerByNameCollecter(region, name))
	}
	summoner, err := db.QuerySummonerByName(region, name)
	if err != nil {
		c.Status(404).SendString(err.Error())
		return err
	}
	c.Status(200).JSON(*summoner)
	return nil
}

// update/summoner/:puuid
func UpdateSummoner(c *fiber.Ctx) error {
	puuid := c.Params("puuid")
	if !db.SummonerPuuidExists(puuid) {
		c.SendStatus(404)
		return errors.New(puuid + " not in db")
	}
	summoner, err := db.QuerySummonerByPuuid(puuid)
	if err != nil {
		c.Status(404).SendString(err.Error())
		return err
	}
	err = <-prioritySummonerScheduler.Schedule(collection.NewSummonerByPuuidCollecter(summoner.Region, summoner.Puuid))
	if err != nil {
		c.Status(500).SendString(err.Error())
		return err
	}
	c.Status(200).Status(200).SendString("updated!")
	return nil
}

// matches/:puuid
func MatchHistory(c *fiber.Ctx) error {
	puuid := c.Params("puuid")
	count := 10
	matchHistory, err := db.QueryMatchHistory(puuid, count)
	if err != nil {
		c.Status(500).SendString(err.Error())
		return err
	}
	c.Status(200).JSON(matchHistory)
	return nil
}

// update/matches/:puuid
func UpdateMatchHistory(c *fiber.Ctx) error {
	puuid := c.Params("puuid")
	err := <-priorityMatchScheduler.Schedule(collection.NewMatchHistoryCollecter(puuid, 0, true))
	if err != nil {
		c.Status(500).SendString(err.Error())
		return err
	}
	c.Status(200).SendString("updated!")
	return nil
}
