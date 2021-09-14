package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/fatih/color"
)

type price_data struct {
	Price                    float64 `json:"price"`
	Volume_24h               float64 `json:"volume_24h"`
	Percent_change_1h        float64 `json:"percent_change_1h"`
	Percent_change_24h       float64 `json:"percent_change_24h"`
	Percent_change_7d        float64 `json:"percent_change_7d"`
	Percent_change_30d       float64 `json:"percent_change_30d"`
	Percent_change_60d       float64 `json:"percent_change_60d"`
	Percent_change_90d       float64 `json:"percent_change_90d"`
	Market_cap               float64 `json:"market_cap"`
	Market_cap_dominance     float64 `json:"market_cap_dominance"`
	Fully_diluted_market_cap float64 `json:"fully_diluted_market_cap"`
	Last_updated             string  `json:"last_updated"`
}

// I'd get rid of this if I knew a way
type usd_data struct {
	USD price_data `json:"USD"`
}

type cryptocurrency struct {
	Id                 string   `json:"id"`
	Name               string   `json:"name"`
	Symbol             string   `json:"symbol"`
	Slug               string   `json:"slug"`
	Num_market_pairs   int      `json:"num_market_pairs"`
	Date_added         string   `json:"date_added"`
	Tags               []string `json:"tags"`
	Max_supply         float64  `json:"max_supply"`
	Circulating_supply float64  `json:"circulating_supply"`
	Total_supply       float64  `json:"total_supply"`
	Platform           string   `json:"platform"`
	Cmc_rank           int      `json:"cmc_rank"`
	Last_updated       string   `json:"last_updated"`
	Quote              usd_data `json:"quote"`
}

func main() {

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	// get coinmarketcap key from env
	key := os.Getenv("CMC_KEY")
	if len(key) == 0 {
		color.Red("Error, coinmarketcap key not set to CMC_KEY in env")
		return 
	}

	q := url.Values{}
	q.Add("start", "1")
	q.Add("limit", "20")
	q.Add("convert", "USD")

	req.Header.Set("Accepts", "Application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", key)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error contacting server")
		log.Print(err)
		os.Exit(1)
	}

	// dig through the json for the data I want
	respBody, _ := ioutil.ReadAll(resp.Body)
	var resJSON map[string]interface{}
	json.Unmarshal([]byte(respBody), &resJSON)
	dataStr, _ := json.Marshal(resJSON["data"])
	var coins []cryptocurrency
	json.Unmarshal(dataStr, &coins)

	// setup colors
	nameText := color.New(color.FgBlack).Add(color.BgCyan)
	upText := color.New(color.FgGreen)
	downText := color.New(color.FgRed)

	// print data
	for _, coin := range coins {
		nameText.Printf("%s: $%v", coin.Name, coin.Quote.USD.Price)

		dayChange := coin.Quote.USD.Percent_change_24h
		weekChange := coin.Quote.USD.Percent_change_7d

		if dayChange > 0 {
			upText.Printf(", 24h: " + fmt.Sprint(dayChange) + "%%, ")
		} else {
			downText.Printf(", 24h: " + fmt.Sprint(dayChange) + "%%, ")
		}

		if weekChange > 0 {
			upText.Printf("7d: " + fmt.Sprint(weekChange) + "%%\n")
		} else {
			downText.Printf("7d: " + fmt.Sprint(weekChange) + "%%\n")
		}
	}

}
