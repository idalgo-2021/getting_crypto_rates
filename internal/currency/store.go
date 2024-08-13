package currency

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

type Currency struct {
	Code  string
	Name  string
	Rates map[string]float64
}

// Хранилище валют
type CurrencyStore struct {
	sync.Mutex
	Currencies map[string]Currency
}

func NewCurrencyStore() *CurrencyStore {
	return &CurrencyStore{
		Currencies: make(map[string]Currency),
	}
}

// Получить список валют
func (cs *CurrencyStore) FetchAllCurrencies() error {
	resp, err := http.Get("https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies.json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	csData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var csMap map[string]string
	err = json.Unmarshal(csData, &csMap)
	if err != nil {
		return err
	}

	for code, name := range csMap {
		cs.Currencies[code] = Currency{
			Code:  code,
			Name:  name,
			Rates: make(map[string]float64),
		}
	}
	return nil
}

func (cs *CurrencyStore) UpdateCurrency(currency Currency) {
	cs.Lock()
	defer cs.Unlock()
	cs.Currencies[currency.Code] = currency
}
