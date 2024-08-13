package currency

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CurrencyFetcher struct{}

func NewCurrencyFetcher() *CurrencyFetcher {
	return &CurrencyFetcher{}
}

func (cf *CurrencyFetcher) FetchCurrencyRates(currencyCode string) (map[string]float64, error) {
	resp, err := http.Get(fmt.Sprintf("https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/%s.json", currencyCode))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ratesData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ratesMap map[string]interface{}
	err = json.Unmarshal(ratesData, &ratesMap)
	if err != nil {
		return nil, err
	}

	currencyRates := make(map[string]float64)
	for code, rate := range ratesMap[currencyCode].(map[string]interface{}) {
		currencyRates[code] = rate.(float64)
	}
	return currencyRates, nil
}
