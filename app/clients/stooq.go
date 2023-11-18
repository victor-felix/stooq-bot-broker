package clients

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/victor-felix/chat-bot/app/pkg/httpclient"
)

type stooqClient struct {
	baseUrl string
	httpClient httpclient.HTTPClient
	log zerolog.Logger
}

func NewStooqClient(baseUrl string, log zerolog.Logger) *stooqClient {
	return NewStooqClientWithHTTPClient(baseUrl, log, &http.Client{
		Timeout: time.Duration(60) * time.Second,
	})
}

func NewStooqClientWithHTTPClient(baseUrl string, log zerolog.Logger, httpClient httpclient.HTTPClient) *stooqClient {
	return &stooqClient{
		baseUrl: baseUrl,
		log: log,
		httpClient: httpClient,
	}
}

func (sc *stooqClient) GetStockPrice(symbol string) (string, error) {
	url := fmt.Sprintf("%s/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv", sc.baseUrl, url.QueryEscape(symbol))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		sc.log.Error().Msg(err.Error())
		return "", err
	}

	res, err := sc.httpClient.Do(req)
	if err != nil {
		sc.log.Error().Msg(err.Error())
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		sc.log.Error().Msg("stooq api returned a non 200 status code")
		return "", err
	}

	csvContent, err := csv.NewReader(res.Body).ReadAll()
	if err != nil {
		sc.log.Error().Msg(err.Error())
		return "", err
	}

	symbol = csvContent[1][0]
	close := csvContent[1][6]

	if close == "N/D" {
		return fmt.Sprintf("%s quote is not available", strings.ToUpper(symbol)), nil
	}

	return fmt.Sprintf("%s quote is $%s per share", strings.ToUpper(symbol), close), nil
}