package bot

import "github.com/adshao/go-binance/v2"

func (b *Bot) GetCurrentCryptoPrice(symbol string) (<-chan string, chan<- struct{}, error) {
	currentPriceC := make(chan string)
	var wsErr error

	eventHandler := func(event *binance.WsMarketStatEvent) {
		currentPriceC <- event.LastPrice
	}

	errHandler := func(err error) {
		wsErr = err
	}

	doneC, stopC, err := binance.WsMarketStatServe(symbol, eventHandler, errHandler)
	if wsErr != nil {
		return nil, nil, wsErr
	}

	if err != nil {
		return nil, nil, err
	}

	go func() {
		for {
			select {
			case <-doneC:
				close(currentPriceC)
				return
			case <-stopC:
				close(currentPriceC)
				return
			}
		}
	}()

	return currentPriceC, stopC, nil
}
