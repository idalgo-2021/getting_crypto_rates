package main

import (
	"fmt"
	"os"
	"time"

	"getting_crypto_rates/internal/currency"
)

func main() {

	// Инициализация хранилища валют и получение всех валют
	currencyStore := currency.NewCurrencyStore()
	err := currencyStore.FetchAllCurrencies()
	if err != nil {
		panic(err)
	}

	// Инициализация каналов
	currencyChan := make(chan currency.Currency, len(currencyStore.Currencies))
	resultChan := make(chan currency.Currency, len(currencyStore.Currencies))

	// Запуск пула из n-горутин
	n := 5
	workerPool := currency.NewWorkerPool(n, currencyChan, resultChan)
	workerPool.Start()

	startTime := time.Now()
	resultCount := 0

	// Отправка валют в канал для обработки
	for _, curr := range currencyStore.Currencies {
		currencyChan <- curr
	}

	// Обработка результатов
	for {
		if resultCount == len(currencyStore.Currencies) {
			close(currencyChan)
			break
		}
		select {
		case c := <-resultChan:
			currencyStore.UpdateCurrency(c)
			resultCount++
		case <-time.After(3 * time.Second):
			fmt.Println("Тайм-аут")
			return
		}
	}

	endTime := time.Now()

	// Вывод результатов в консоль
	fmt.Println("======== Результаты ========")
	for _, curr := range currencyStore.Currencies {
		fmt.Printf("%s (%s): %d курсов\n", curr.Name, curr.Code, len(curr.Rates))
	}
	fmt.Println("============================")
	fmt.Println("Время выполнения: ", endTime.Sub(startTime))

	//Вывод результатов в файл
	file, err := os.Create("currency_rates.txt")
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer file.Close()
	file.WriteString("======== Результаты ========\n")
	for _, curr := range currencyStore.Currencies {
		file.WriteString(fmt.Sprintf("Валюта: %s (%s)\n", curr.Name, curr.Code))
		for otherCode, rate := range curr.Rates {
			line := fmt.Sprintf("  1 %s = %.4f %s\n", curr.Code, rate, otherCode)
			file.WriteString(line)
		}
		file.WriteString("\n")
	}
	file.WriteString("============================\n")
	file.WriteString("Время выполнения: " + endTime.Sub(startTime).String() + "\n")
}
