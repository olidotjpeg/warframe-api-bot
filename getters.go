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

	var responseObject VoidTrader
	json.Unmarshal(responseData, &responseObject)

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
