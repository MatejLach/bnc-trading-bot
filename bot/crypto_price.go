package bot

import (
	"time"

	"github.com/adshao/go-binance/v2"
)

func (b *Bot) GetCurrentCryptoPrice(symbol string) (<-chan string, chan<- struct{}, error) {
	currentPriceC := make(chan string)
	stopC := make(chan struct{})
	var done bool
	var wsErr error

	eventHandler := func(event *binance.WsMarketStatEvent) {
		if !done {
			currentPriceC <- event.LastPrice
		}
	}

	errHandler := func(err error) {
		wsErr = err
	}

	doneC, stopWsC, err := binance.WsMarketStatServe(symbol, eventHandler, errHandler)
	if wsErr != nil {
		return nil, nil, wsErr
	}

	if err != nil {
		return nil, nil, err
	}

	go func() {
		defer close(currentPriceC)

		for {
			select {
			case <-doneC:
				time.Sleep(5 * time.Second)
				done = true
				return
			case <-stopC:
				stopWsC <- struct{}{}
				time.Sleep(5 * time.Second)
				done = true
				return
			}
		}
	}()

	return currentPriceC, stopC, nil
}
