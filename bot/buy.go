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

func (b *Bot) BuyIfDecreaseByPercent(symbolPriceChan <-chan string, symbolPriceCloseChan chan<- struct{}, buyCfg BuyConfig) error {
	bimOriginalPrice, err := money.ParseBimoney(<-symbolPriceChan)
	if err != nil {
		return err
	}

	configuredPercentDiff, err := money.ParseBimoney(strconv.Itoa(buyCfg.PercentageDiff))
	if err != nil {
		return err
	}

	if buyCfg.TargetPriceToBuyAt == "" || buyCfg.TargetPriceToBuyAt == "0" {
		if configuredPercentDiff == 0 {
			return fmt.Errorf("config: either 'target_price_to_buy_at' or 'percentage_diff' have to be configured for %s%s",
				buyCfg.BuySymbol, buyCfg.BuyWithHoldingSymbol)
		}

		buyCfg.TargetPriceToBuyAt = (bimOriginalPrice + bimOriginalPrice.AmountFromPercentage(configuredPercentDiff)).FormatBimoney(true)
	} else {
		bimTargetToBuyAt, err := money.ParseBimoney(buyCfg.TargetPriceToBuyAt)
		if err != nil {
			return err
		}

		if bimTargetToBuyAt >= bimOriginalPrice {
			return fmt.Errorf("config: 'target_price_to_buy_at' for %s%s has to be lower than the current price of the symbol", buyCfg.BuySymbol, buyCfg.BuyWithHoldingSymbol)
		}

		configuredPercentDiff = bimOriginalPrice.PercentageChange(bimTargetToBuyAt)
	}

	log.Printf("Starting with the initial %s price of %s %s, with a purchase target of %s or less, (price decrease of at least %s percentage points), before buying...",
		buyCfg.BuySymbol, bimOriginalPrice.FormatBimoney(false), buyCfg.BuyWithHoldingSymbol, buyCfg.TargetPriceToBuyAt, configuredPercentDiff.FormatBimoney(false))

	// enforce Binance server time
	_, err = b.binanceClient.NewSetServerTimeService().Do(context.Background())
	if err != nil {
		return err
	}

	for {
		currentPrice := <-symbolPriceChan

		bimCurrentPrice, err := money.ParseBimoney(currentPrice)
		if err != nil {
			return err
		}

		currentPercentDiff := bimOriginalPrice.PercentageChange(bimCurrentPrice)

		if currentPercentDiff <= configuredPercentDiff && bimCurrentPrice < bimOriginalPrice {
			log.Printf("Price decreased from %s to %s, which is a %s percent decrease!",
				bimOriginalPrice.FormatBimoney(false), currentPrice, currentPercentDiff.FormatBimoney(false))

			cryptoQtyToBuy, err := money.ParseBimoney(buyCfg.BuyQuantity)
			if err != nil {
				cryptoQtyToBuy = 0
			}

			if buyCfg.BuyForPercentOfHoldings != 0 {
				bimBuyPercentOfHoldings, err := money.ParseBimoney(strconv.Itoa(int(buyCfg.BuyForPercentOfHoldings)))
				if err != nil {
					return err
				}

				balance, err := b.GetAccountBalance(buyCfg.BuyWithHoldingSymbol)
				if err != nil {
					return err
				}

				cryptoQtyToBuy = balance.AmountFromPercentage(bimBuyPercentOfHoldings).PortionOf(bimCurrentPrice)
			}

			err = b.buyCrypto(cryptoQtyToBuy, buyCfg, currentPrice)
			if err != nil {
				if errors.Is(err, InsufficientFunds) {
					log.Print(fmt.Errorf("%w, retrying in 1min", err).Error())
					time.Sleep(1 * time.Minute)
					continue
				}
				return err
			}

			symbolPriceCloseChan <- struct{}{}
			break
		}
	}

	return nil
}

func (b *Bot) buyCrypto(cryptoQuantity money.Bimoney, buyCfg BuyConfig, currentPrice string) error {
	strQuantity := cryptoQuantity.FormatBimoney(true)

	log.Printf("Buying %s for %s...", buyCfg.BuySymbol, buyCfg.BuyWithHoldingSymbol)

	orderResp, err := b.binanceClient.NewCreateOrderService().
		Symbol(fmt.Sprintf("%s%s", buyCfg.BuySymbol, buyCfg.BuyWithHoldingSymbol)).
		Side(binance.SideTypeBuy).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(strQuantity).
		Price(currentPrice).Do(context.Background())

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "insufficient balance") {
			return InsufficientFunds
		}
		return err
	}

	if orderResp.Status == binance.OrderStatusTypeFilled && len(orderResp.Fills) > 0 {
		log.Printf("Purchased %s of %s for %s", orderResp.Fills[0].Quantity, buyCfg.BuySymbol, orderResp.Fills[0].Price)
	}

	return nil
}
