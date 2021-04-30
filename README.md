bnc-trading-bot
=


[![License: GPL v3](https://img.shields.io/badge/GPLv3-Free%20as%20in%20freedom-blue)](LICENSE)

> [Download the latest release for your OS](https://github.com/MatejLach/bnc-trading-bot/releases)

This bot automatically trades crypto / fiat assets on [Binance.com](https://www.binance.com/en) according to simple, configurable rules.
It allows you to buy crypto assets when they go down in value by a configurable percentage amount and sell when they're up, automatically.

Building from source
-
Assuming a correctly set up [Go toolchain](https://golang.org/doc/install), simply run:

`$ go build` 

from within the cloned source directory. 

Running
-
Navigate with your shell to the directory where the compiled binary is located and execute with:

`$ ./bnc-trading-bot`

Configuration
-

Upon first launch, a new `config.json` configuration file will be created that looks something like the following, (excluding the explanatory comments):

```
{
    "api_key": "<YOUR-BINANCE-API-KEY>", // replace everything between " " with your personal Binance.com API key
    "api_secret": "<YOUR-BINANCE-API-SECRET>", // replace everything between " " with your personal Binance.com API secret
    "Sell": [
        {
            "sell_holding_symbol": "BNB", // i.e. "BNB" as in 'I have BNB holdings'
            "sell_for_symbol": "GBP", // that I want to sell for 'GBP'
            "target_price_to_sell_at": "524.00", // sell when "BNB" price against "GBP" reaches an exact amount, set 'percentage_diff' to 0 for this to take effect (has to be higher than the current price)
            "percentage_diff": 55, // OR when 'BNB' is up 55% against 'GBP' from when I started this program (has to be positive)
            "sell_percent_of_holdings": 25, // sell 25% of my total BNB wallet balance, set to 0 in order for the 'sell_quantity' setting to take effect
            "sell_quantity": "0.45" // OR sell exactly 0.45 BNB when it's up 55%, 'sell_percent_of_holdings' takes precedence if it's set to > 0
        }
    ],
    "Buy": [
        {
            "buy_symbol": "BNB", // purchase 'BNB' 
            "buy_with_holding_symbol": "GBP", // for money from my 'GBP' wallet
            "target_price_to_buy_at": "192", // buy when "BNB" price against "GBP" reaches an exact amount, set 'percentage_diff' to 0 for this to take effect (has to be lower than the current price)
            "percentage_diff": -15, // OR once 'BNB' is 15% down against 'GBP' from when I started this program (has to be negative)
            "buy_for_percent_of_holdings": 85, // buy as much BNB as possible for 85% of your GBP fiat wallet balance
            "buy_quantity": "1.12" // OR buy exactly 1.12 BNB if 'buy_for_percent_of_holdings' is set to 0, otherwise 'buy_for_percent_of_holdings' takes precedence
        }
    ]
}
```

You can specify as many Buy/Sell configurations as you want by placing successive configuration objects `{...}` separated by `,` in between `[]`, see [this](https://opensource.adobe.com/Spry/samples/data_region/JSONDataSetSample.html#Example2) if you are unfamiliar with JSON arrays.

The individual buy/sell configurations will run in parallel.

To obtain your personal Binance API key/secret pair, consult [the relevant support article](https://www.binance.com/en/support/articles/360002502072).
Make sure to **never share your config.json with anyone if it has your API details filled in**, otherwise they will be able to do trades on your behalf.

Contributing
-

Bug reports and pull requests are welcome. Do not hesitate to open a PR / file an issue, or a feature request.

Disclaimer
-

This program is provided 'AS IS', without any warranty or assumed liability by its creator(s).
Any financial looses due its usage are solely the responsibility of its user(s).

**Use at your own risk!**
