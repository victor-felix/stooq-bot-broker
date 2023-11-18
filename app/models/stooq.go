package models

type StooqClient interface {
	GetStockPrice(stockName string) (string, error)
}