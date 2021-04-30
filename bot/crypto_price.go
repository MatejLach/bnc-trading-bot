package bot

import (
	"github.com/adshao/go-binance/v2"
	"time"
)

func (b *Bot) GetCurrentCryptoPrice(symbol string) (string, error) {
	var wsErr error
	var currentPrice string

	eventHandler := func(event *binance.WsMarketStatEvent) {
		currentPrice = event.LastPrice
	}

	errHandler := func(err error) {
		wsErr = err
	}

	_, stopC, err := binance.WsMarketStatServe(symbol, eventHandler, errHandler)
	if wsErr != nil {
		return "", wsErr
	}

	if err != nil {
		return "", err
	}

	for {
		if currentPrice != "" {
			stopC <- struct{}{}
			return currentPrice, nil
		} else {
			time.Sleep(1 * time.Second)
			continue
		}
	}
}
