package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
)

type Config struct {
	APIKey    string       `json:"api_key"`
	APISecret string       `json:"api_secret"`
	Sell      []SellConfig `json:"Sell"`
	Buy       []BuyConfig  `json:"Buy"`
}

type SellConfig struct {
	SellHoldingSymbol      string `json:"sell_holding_symbol"`
	SellForSymbol          string `json:"sell_for_symbol"`
	TargetPriceToSellAt    string `json:"target_price_to_sell_at"`
	PercentageDiff         uint   `json:"percentage_diff"`
	SellPercentOfHoldings  uint   `json:"sell_percent_of_holdings"`
	SellQuantityOfHoldings string `json:"sell_quantity"`
}

type BuyConfig struct {
	BuySymbol               string `json:"buy_symbol"`
	BuyWithHoldingSymbol    string `json:"buy_with_holding_symbol"`
	TargetPriceToBuyAt      string `json:"target_price_to_buy_at"`
	PercentageDiff          int    `json:"percentage_diff"`
	BuyForPercentOfHoldings uint   `json:"buy_for_percent_of_holdings"`
	BuyQuantity             string `json:"buy_quantity"`
}

func initConfig() {
	log.Println("No config.json found, attempting to create a new config...")
	config := Config{
		Sell: []SellConfig{
			{
				SellHoldingSymbol:      "",
				SellForSymbol:          "",
				PercentageDiff:         0,
				SellPercentOfHoldings:  0,
				SellQuantityOfHoldings: "",
			},
		},
		Buy: []BuyConfig{
			{
				BuySymbol:               "",
				BuyWithHoldingSymbol:    "",
				PercentageDiff:          0,
				BuyForPercentOfHoldings: 0,
				BuyQuantity:             "",
			},
		},
	}

	jsonBytes, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("config.json", jsonBytes, 0777)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("New config.json created, come back once you're filled it with your desired parameters.")
}

func parseConfig() (Config, error) {
	cfgData, err := os.ReadFile("config.json")
	if err != nil {
		return Config{}, err
	}

	var config Config

	err = json.Unmarshal(cfgData, &config)
	if err != nil {
		return Config{}, err
	}

	if config.APIKey == "" || config.APISecret == "" {
		return Config{}, fmt.Errorf("you need to provide Binance API key and API secret in the config to be able to use this bot")
	}

	if len(config.Sell) == 0 && len(config.Buy) == 0 {
		return Config{}, fmt.Errorf("you need to fill in at least one 'Sell' or 'Buy' section of the config to be able to use this bot")
	}

	for _, buyCfg := range config.Buy {
		if buyCfg.PercentageDiff != 0 && !math.Signbit(float64(buyCfg.PercentageDiff)) {
			return Config{}, fmt.Errorf("buy config: asset price percentage decrease for %s needs to be negative (-%d) to be considered valid", buyCfg.BuySymbol, buyCfg.PercentageDiff)
		}
	}

	return config, nil
}

func configExists() bool {
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		return false
	}

	return true
}
