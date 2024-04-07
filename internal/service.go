package internal

import (
	"binanceAPI/pkg/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	exchangeInfoEndpoint    = "https://api.binance.com/api/v3/exchangeInfo"
	priceTickerEndpointTmpl = "https://api.binance.com/api/v3/ticker/price?symbol=%s"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetExchangeInfo() (ExchangeInfo, error) {
	response, err := http.Get(exchangeInfoEndpoint)
	if err != nil {
		logger.Errorf("failed retrieve exchange info: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		logger.Errorf("unexpected response status on exchange info: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Errorf("failed read exchange info response: error:%v", err)
	}

	var data ExchangeInfo
	if err = json.Unmarshal(body, &data); err != nil {
		logger.Errorf("failed unmarshal exchange info response: error:%v, body:%s", err, string(body))
	}

	return data, nil
}

func (s *Service) GetSymbolPrice(symbol string) (SymbolPrice, error) {
	response, err := http.Get(fmt.Sprintf(priceTickerEndpointTmpl, symbol))
	if err != nil {
		logger.Errorf("failed retrieve symbol price: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		logger.Errorf("unexpected response status on symbol price: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Errorf("failed read symbol price response: error:%v", err)
	}

	var data SymbolPrice
	if err = json.Unmarshal(body, &data); err != nil {
		logger.Errorf("failed unmarshal symbol price response: error:%v, body:%s", err, string(body))
	}

	return data, nil
}
