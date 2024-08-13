package currency

import "fmt"

type WorkerPool struct {
	workers      int
	currencyChan <-chan Currency
	resultChan   chan<- Currency
	fetcher      *CurrencyFetcher
}

func NewWorkerPool(workers int, currencyChan <-chan Currency, resultChan chan<- Currency) *WorkerPool {
	return &WorkerPool{
		workers:      workers,
		currencyChan: currencyChan,
		resultChan:   resultChan,
		fetcher:      NewCurrencyFetcher(),
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		go wp.runWorker(i)
	}
}

func (wp *WorkerPool) runWorker(workerId int) {
	fmt.Printf("Worker %d запущен\n", workerId)
	for currency := range wp.currencyChan {
		rates, err := wp.fetcher.FetchCurrencyRates(currency.Code)
		if err != nil {
			fmt.Printf("Ошибка при получении курсов валют для %s: %v\n", currency.Code, err)
			continue
		}
		currency.Rates = rates
		wp.resultChan <- currency
	}
	fmt.Printf("Worker %d остановлен\n", workerId)
}
