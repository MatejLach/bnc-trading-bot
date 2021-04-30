package bot

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/rs/zerolog/log"

	"github.com/MatejLach/bnc-trading-bot/money"
)

func (b *Bot) SellIfIncreaseByPercent(originalPrice string, sellCfg SellConfig) error {
	bimOriginalPrice, err := money.ParseBimoney(originalPrice)
	if err != nil {
		return err
	}

	configuredPercentDiff, err := money.ParseBimoney(strconv.Itoa(int(sellCfg.PercentageDiff)))
	if err != nil {
		return err
	}

	if sellCfg.TargetPriceToSellAt == "" || sellCfg.TargetPriceToSellAt == "0" {
		if configuredPercentDiff == 0 {
			return fmt.Errorf("config: either 'target_price_to_sell_at' or 'percentage_diff' have to be configured for %s%s",
				sellCfg.SellHoldingSymbol, sellCfg.SellForSymbol)
		}

		sellCfg.TargetPriceToSellAt = (bimOriginalPrice + bimOriginalPrice.AmountFromPercentage(configuredPercentDiff)).FormatBimoney(true)
	} else {
		bimTargetToSellAt, err := money.ParseBimoney(sellCfg.TargetPriceToSellAt)
		if err != nil {
			return err
		}

		if bimTargetToSellAt <= bimOriginalPrice {
			return fmt.Errorf("config: 'target_price_to_sell_at' for %s%s has to be greater than the current price of the symbol", sellCfg.SellHoldingSymbol, sellCfg.SellForSymbol)
		}

		configuredPercentDiff = bimOriginalPrice.PercentageChange(bimTargetToSellAt)
	}

	log.Printf("Starting with the initial %s price of %s %s, with a sell target of %s or more, (price increase of at least %s percentage points), before selling...",
		sellCfg.SellHoldingSymbol, originalPrice, sellCfg.SellForSymbol, sellCfg.TargetPriceToSellAt, configuredPercentDiff.FormatBimoney(false))

	// enforce Binance server time
	_, err = b.binanceClient.NewSetServerTimeService().Do(context.Background())
	if err != nil {
		return err
	}

	for {
		currentPrice, err := b.GetCurrentCryptoPrice(fmt.Sprintf("%s%s", sellCfg.SellHoldingSymbol, sellCfg.SellForSymbol))
		if err != nil {
			log.Err(err).Msg("sleeping for 5s to try & recover")
			time.Sleep(5)
			continue
		}

		bimCurrentPrice, err := money.ParseBimoney(currentPrice)
		if err != nil {
			return err
		}

		currentPercentDiff := bimOriginalPrice.PercentageChange(bimCurrentPrice)

		if currentPercentDiff >= configuredPercentDiff && bimCurrentPrice > bimOriginalPrice && bimCurrentPrice != bimOriginalPrice {
			log.Printf("Price increased from %s to %s, which is a %s percent increase!",
				originalPrice, currentPrice, currentPercentDiff.FormatBimoney(false))

			cryptoQtyToSell, err := money.ParseBimoney(sellCfg.SellQuantityOfHoldings)
			if err != nil {
				return err
			}

			if sellCfg.SellPercentOfHoldings != 0 {
				bimSellPercentOfHoldings, err := money.ParseBimoney(strconv.Itoa(int(sellCfg.SellPercentOfHoldings)))
				if err != nil {
					return err
				}

				balance, err := b.GetAccountBalance(sellCfg.SellHoldingSymbol)
				if err != nil {
					return err
				}

				cryptoQtyToSell = balance.AmountFromPercentage(bimSellPercentOfHoldings)
			}

			err = b.sellCrypto(cryptoQtyToSell, sellCfg, currentPrice)
			if err != nil {
				if errors.Is(err, InsufficientFunds) {
					log.Print(fmt.Errorf("%w, retrying in 1min", err).Error())
					time.Sleep(1 * time.Minute)
					continue
				}
				return err
			}

			// TODO: Do we want to continue if there's sufficient funds left to satisfy sellCfg conditions again?
			bimOriginalPrice = bimCurrentPrice
			break
		}
	}

	return nil
}

func (b *Bot) sellCrypto(cryptoQuantity money.Bimoney, sellCfg SellConfig, currentPrice string) error {
	strQuantity := cryptoQuantity.FormatBimoney(true)

	log.Printf("Selling %s for %s...", sellCfg.SellHoldingSymbol, sellCfg.SellForSymbol)

	orderResp, err := b.binanceClient.NewCreateOrderService().
		Symbol(fmt.Sprintf("%s%s", sellCfg.SellHoldingSymbol, sellCfg.SellForSymbol)).
		Side(binance.SideTypeSell).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(strQuantity).
		Price(currentPrice).Do(context.Background())

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "insufficient balance") {
			return InsufficientFunds
		}
		return err
	}

	if orderResp.Status == binance.OrderStatusTypeFilled && len(orderResp.Fills) > 0 {
		log.Printf("Sold %s of %s for %s", orderResp.Fills[0].Quantity, sellCfg.SellHoldingSymbol, orderResp.Fills[0].Price)
	}

	return nil
}
