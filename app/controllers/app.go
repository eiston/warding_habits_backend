package controllers

import (
	"github.com/revel/revel"
	"net/http"
	"log"
	"fmt"
	"encoding/json"
)

type App struct {
	*revel.Controller
}

type Point struct {
	x int
	y int
	value int
}

const API_KEY = "RGAPI-ca13e0ef-8468-4771-a54d-40881a73496b"

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) SummonerInfo(region string, name string) revel.Result {
	var url	=  "https://na1.api.riotgames.com"
	switch region {
	case "na1":
		url = "https://na1.api.riotgames.com"
	default:
		url	=  "https://na1.api.riotgames.com"
	}
	url += "/lol/summoner/v4/summoners/by-name/"
	url += name
	url += "?api_key=RGAPI-ca13e0ef-8468-4771-a54d-40881a73496b"


	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	
	defer resp.Body.Close()
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	data := make(map[string]interface{})
	data["error"] = nil
	data["data"] = result
	fmt.Printf("%s\n",data)
	return c.RenderJSON(data)
}

func (c App) getSummonerByName(region string, name string) map[string]interface{} {
	var url = fmt.Sprintf("https://%s.api.riotgames.com/lol/summoner/v4/summoners/by-name/%s?api_key=%s", region, name, API_KEY)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result
}

func (c App) getMatchlist(region string, accountId string, beginTime int) map[string]interface{} {
	var url = fmt.Sprintf("https://%s.api.riotgames.com/lol/summoner/v4/summoners/by-name/%s?beginTime=%d&api_key=%s", region, accountId, beginTime, API_KEY)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result
}

func (c App) WardData(region string, name string) revel.Result {
	result := c.getSummonerByName(region, name)

	accountId := result["data"]["accountId"]

	matchList := c.getMatchlist(region, accountId, 0)

	for _, match := range matchList {
		
	}




	var data [max_length]([]Point)
}