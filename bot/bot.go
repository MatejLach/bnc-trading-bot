package bot

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/hashicorp/go-retryablehttp"
	_ "github.com/mattn/go-sqlite3"
)

var (
	InsufficientFunds = errors.New("account has insufficient balance for requested action")
)

type Bot struct {
	binanceClient *binance.Client
	db            *sql.DB
	Config        Config
}

func New() *Bot {
	if !configExists() {
		initConfig()
		os.Exit(0)
	}

	config, err := parseConfig()
	if err != nil {
		log.Fatal(err)
	}

	binance.WebsocketKeepalive = true

	client := binance.NewClient(config.APIKey, config.APISecret)

	httpClient := retryablehttp.NewClient()
	httpClient.Logger = nil // disable default DEBUG logs

	httpClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if errors.Is(err, syscall.ECONNRESET) {
			return true, err
		}
		return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
	}
	httpClient.RetryWaitMin = 30 * time.Second
	httpClient.RetryWaitMax = 2 * time.Minute

	client.HTTPClient = httpClient.StandardClient()

	// TODO: Add on-disk state keeping
	sqliteDb, err := sql.Open("sqlite3", "./state.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	//err =initDb(sqliteDb)
	//if err != nil {
	//	log.Fatalln(err)
	//}

	return &Bot{binanceClient: client, db: sqliteDb, Config: config}
}
