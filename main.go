package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/MatejLach/bnc-trading-bot/bot"
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func main() {
	botClient := bot.New()
	wg := sync.WaitGroup{}

	for _, cfg := range botClient.Config.Sell {
		startingPrice, err := botClient.GetCurrentCryptoPrice(fmt.Sprintf("%s%s", cfg.SellHoldingSymbol, cfg.SellForSymbol))
		if err != nil {
			log.Warn().Err(err).Msgf("failed getting initial price of %s, skipping...", cfg.SellHoldingSymbol)
			continue
		}

		wg.Add(1)

		go func(sellCfg bot.SellConfig) {
			defer wg.Done()

			err := botClient.SellIfIncreaseByPercent(startingPrice, sellCfg)
			if err != nil {
				symbol := fmt.Sprintf("%s%s", sellCfg.SellHoldingSymbol, sellCfg.SellForSymbol)
				log.Warn().Err(err).Msgf("failed executing sell order for '%s', skipping...", symbol)
			}

		}(cfg)
	}

	for _, cfg := range botClient.Config.Buy {
		startingPrice, err := botClient.GetCurrentCryptoPrice(fmt.Sprintf("%s%s", cfg.BuySymbol, cfg.BuyWithHoldingSymbol))
		if err != nil {
			log.Warn().Err(err).Msgf("failed getting initial price of %s, skipping...", cfg.BuySymbol)
			continue
		}

		wg.Add(1)

		go func(buyCfg bot.BuyConfig) {
			defer wg.Done()

			err := botClient.BuyIfDecreaseByPercent(startingPrice, buyCfg)
			if err != nil {
				symbol := fmt.Sprintf("%s%s", buyCfg.BuySymbol, buyCfg.BuyWithHoldingSymbol)
				log.Warn().Err(err).Msgf("failed executing purchase order for '%s', skipping...", symbol)
			}

		}(cfg)
	}

	wg.Wait()
}
