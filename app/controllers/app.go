package controllers

import (
	"github.com/revel/revel"
	"net/http"
	"log"
	"fmt"
	"encoding/json"
	// "sort"
)

type App struct {
	*revel.Controller
}

type Point struct {
	timestamp int
	x int
	y int
	value int
}


const API_KEY = "RGAPI-bad956ea-222a-490a-abb9-4aa9e73c3f88"

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
	var url = fmt.Sprintf("https://%s.api.riotgames.com/lol/summoner/v4/matchlists/by-account/%s?beginTime=%d&api_key=%s", region, accountId, beginTime, API_KEY)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result
}
func (c App) getMatchDetails(region string, gameId string) map[string]interface{} {
	var url = fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v4/matches/%s?beginTime=%d&api_key=%s", region, gameId, API_KEY)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result
}
func (c App) getMatchTimeline(region string, gameId string) map[string]interface{} {
	var url = fmt.Sprintf("https://%s.api.riotgames.com/lol/timelines/by-match/%s?beginTime=%d&api_key=%s", region, gameId, API_KEY)
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
	summoner := c.getSummonerByName(region, name)

	accountId := fmt.Sprintf("%v", summoner["accountId"])

	matchList := c.getMatchlist(region, accountId, 0)

	var data [][]Point
	//Here I dont know how to iterate this shit :(
	for _, match := range list(matchList["matches"]) {
		gameId := fmt.Sprintf("%v", match["gameId"])
		matchDetails := c.getMatchDetails(region, gameId)
		participantId := 1 
		for _, participant := range matchDetails["participantIdentities"] {
			if participant["player"]["accountId"] == summoner["accountId"] {
				participantId = participant["participantId"].(int)
				break
			}
		}

		timeline:= c.getMatchTimeline(region, gameId)

		for _, frame := range timeline["frames"] {
			point := Point{0,0,0,1}
			for _, participant := range frame["participantFrames"] {
				if participantId == participant["participantId"] {
					point.x = participant["position"]["x"].(int)
					point.y = participant["position"]["y"].(int)
				}
			}
			for _, event := range frame["events"] {
				if event["type"] == "WARD_PLACED" && event["creatorId"] == participantId {
					point.timestamp = event["timestamp"].(int)
					if cap(data) <= _ {
						var s []Point
						data = append(data,s) 
					}
					data[_] = append(data[_],point)
				}
			}
		}
	}
	return c.RenderJSON(data)
}