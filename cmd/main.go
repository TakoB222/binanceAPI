package main

import (
	"binanceAPI/internal"
	"binanceAPI/pkg/logger"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

const pairCount = 5

func init() {
	logrus.SetFormatter(new(logrus.TextFormatter))
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	var (
		mu sync.Mutex
		wg sync.WaitGroup
	)

	service := internal.NewService()
	exchangeInfo, err := service.GetExchangeInfo()
	if err != nil {
		logrus.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	messages := make(chan map[string]string)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				close(messages)
				return
			case exchangeRate := <-messages:
				for k, v := range exchangeRate {
					fmt.Printf("%s %s\n", k, v)
				}
			}
		}
	}()

	for i := 0; i < pairCount && len(exchangeInfo.Symbols) >= pairCount; i++ {
		wg.Add(1)
		go func(symbol string) {
			defer func() {
				if r := recover(); r != nil {
					return
				}
			}()
			defer wg.Done()

			ticker := time.NewTicker(500 * time.Millisecond)
			rate := make(map[string]string)
			for {
				select {
				case <-ctx.Done():
					ticker.Stop()
					return
				case <-ticker.C:
					mu.Lock()
					symbolPrice, err := service.GetSymbolPrice(symbol)
					if err != nil {
						logger.Error(err)
					}
					mu.Unlock()
					rate[symbolPrice.Symbol] = symbolPrice.Price
					messages <- rate
				}
			}
		}(exchangeInfo.Symbols[i].Symbol)
	}

	time.Sleep(10 * time.Second)
	cancel()
	wg.Wait()
}
