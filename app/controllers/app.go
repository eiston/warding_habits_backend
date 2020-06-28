package controllers

import (
	"github.com/revel/revel"
	"net/http"
	"log"
	"io/ioutil"
)

type App struct {
	*revel.Controller
}

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

	body, err := ioutil.ReadAll(resp.Body)

	data := make(map[string]interface{})
	data["error"] = nil
	data["data"] = body

	return c.RenderJSON(data)
}