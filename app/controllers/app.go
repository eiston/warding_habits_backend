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

type Summoner struct {
	AccountId string `json:"accountId"`
	ProfileIconId int `json:"profileIconId"`
	RevisionDate float64 `json:"revisionDate"`
	Name string `json:"name"`
	Id string `json:"id"`
	Puuid string `json:"puuid"`
	SummonerLevel float64 `json:"summonerLevel"`
}

type MatchList struct {
	StartIndex int `json:"startIndex"`
	TotalGames int `json:"totalGames"`
	EndIndex int `json:"endIndex"`
	Matches []MatchReference `json:"matches"`
}

type MatchReference struct {
	GameId float64 `json:"gameId"`
	Role string `json:"role"`
	Season int `json:"season"`
	PlatformId string `json:"platformId"`
	Champion int `json:"champion"`
	Queue int `json:"queue"`
	Lane string `json:"lane"`
	Timestamp float64 `json:"timestamp"`
}

type Match struct {
	GameId float64 `json:"gameId"`
	ParticipantIdentities []ParticipantIdentity `json:"participantIdentities"`
	QueueId int `json:"queueId"`
	GameType string `json:"gameType"`
	GameDuration float64 `json:"gameDuration"`
	Teams []TeamStats `json:"teams"`
	PlatformId string `json:"platformId"`
	GameCreation float64 `json:"gameCreation"`
	SeasonId int `json:"seasonId"`
	GameVersion string `json:"gameVersion"`
	MapId int `json:"mapId"`
	GameMode string `json:"gameMode"`
	Participants []Participants `json:"participants"`
}

const API_KEY = "RGAPI-49e44087-8364-429f-9030-ee0df6038eba"

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) SummonerInfo(region string, name string) revel.Result {
	data := c.getSummonerByName(region, name)	

	fmt.Printf("%s\n",data)
	return c.RenderJSON(data)
}

func (c App) getSummonerByName(region string, name string) Summoner {
	var url = fmt.Sprintf("https://%s.api.riotgames.com/lol/summoner/v4/summoners/by-name/%s?api_key=%s", region, name, API_KEY)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(err)
	}
	var result Summoner
	jsonErr := json.Unmarshal(body, &result)
	if jsonErr != nil {
		log.Fatal(err)
	}
	return result
}

func (c App) getMatchlist(region string, accountId string, beginTime int) MatchList {
	var url = fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v4/matchlists/by-account/%s?beginTime=%d&api_key=%s", region, accountId, beginTime, API_KEY)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(err)
	}
	var result MatchList
	jsonErr := json.Unmarshal(body, &result)
	if jsonErr != nil {
		log.Fatal(err)
	}
	return result
}
func (c App) getMatchDetails(region string, gameId string) map[string]interface{} {
	var url = fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v4/matches/%s?api_key=%s", region, gameId, API_KEY)
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
	var url = fmt.Sprintf("https://%s.api.riotgames.com/lol/timelines/by-match/%s?api_key=%s", region, gameId, API_KEY)
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

	b1, _ :=  json.MarshalIndent(summoner, "", " ")
	fmt.Println(string(b1))

	accountId := fmt.Sprintf("%v", summoner.AccountId)

	matchList := c.getMatchlist(region, accountId, 0)

	var data [][]Point
	b2, _ :=  json.MarshalIndent(matchList["matches"], "", " ")
	fmt.Println(string(b2))
	// Here I dont know how to iterate this shit :(
	for _, match := range matchList["matches"] {
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