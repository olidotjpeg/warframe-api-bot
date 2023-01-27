package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func getSortieState() SortieState {
	response, err := http.Get("https://api.warframestat.us/pc/sortie?language=en")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject SortieState
	json.Unmarshal(responseData, &responseObject)

	return responseObject
}

func getVoidTraderState() VoidTrader {
	response, err := http.Get("https://api.warframestat.us/pc/voidTrader/")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	responseObject := VoidTrader{}
	_ = json.Unmarshal(responseData, &responseObject)
	file, err := json.Marshal(responseObject)

	ioutil.WriteFile("voidInventory.json", file, 0644)

	return responseObject
}

func getArbitration() Arbitration {
	response, err := http.Get("https://api.warframestat.us/pc/arbitration/")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject Arbitration
	json.Unmarshal(responseData, &responseObject)

	return responseObject
}

func getWarframeDropData(query string) []WarframeData {
	response, err := http.Get("https://api.warframestat.us/drops/search/" + query)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject []WarframeData
	json.Unmarshal(responseData, &responseObject)

	return responseObject
}

func getItemDetailed(query string) WarframeItem {
	response, err := http.Get("https://api.warframestat.us/items/" + query)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject WarframeItem
	json.Unmarshal(responseData, &responseObject)

	return responseObject
}

func getVoidItemsDetailed() []WarframeItem {
	content, _ := ioutil.ReadFile("voidInventory.json")
	activeInventory := []WarframeItem{}
	mods := []WarframeItem{}

	voidInventory := VoidTrader{}
	json.Unmarshal(content, &voidInventory)

	for _, item := range voidInventory.Inventory {
		newItem := getItemDetailed(item.Item)
		if newItem.Code != 404 {
			activeInventory = append(activeInventory, newItem)
		}

	}

	for _, newItem := range activeInventory {
		if newItem.Category == "Mods" {
			mods = append(mods, newItem)
		}
	}

	return mods
}
