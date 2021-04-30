package bot

import (
	"context"
	"strings"

	"github.com/adshao/go-binance/v2"

	"github.com/MatejLach/bnc-trading-bot/money"
)

func (b *Bot) GetAccountBalance(holdingSymbol string) (money.Bimoney, error) {
	_, err := b.binanceClient.NewSetServerTimeService().Do(context.Background(), binance.WithRecvWindow(0))
	if err != nil {
		return 0, err
	}

	res, err := b.binanceClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		return 0, err
	}

	if len(res.Balances) == 0 {
		return 0, nil
	}

	for _, balance := range res.Balances {
		if strings.ToUpper(holdingSymbol) == strings.ToUpper(balance.Asset) {
			return money.ParseBimoney(balance.Free)
		}
	}

	return 0, nil
}
