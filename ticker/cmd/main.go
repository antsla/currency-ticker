package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/antsla/ticker/internal"
	"github.com/robfig/cron/v3"
)

func main() {
	sources := []internal.HttpStreamSubscriber{
		internal.HttpStreamSubscriber{Url: "http://localhost:8082"},
		internal.HttpStreamSubscriber{Url: "http://localhost:8083"},
	}

	exchgRates := map[string][]TickerHandled{}
	mu := sync.Mutex{}
	wg := &sync.WaitGroup{}

	for _, source := range sources {
		wg.Add(1)
		go handleTickers(&mu, wg, source, exchgRates)
	}

	fmt.Println("Timestamp, IndexPrice")

	c := cron.New()
	_, err := c.AddFunc("* * * * *", func() {
		handleJob(&mu, exchgRates)
	})
	if err != nil {
		log.Printf("cron hasn't been started [%s]", err)
		os.Exit(1)
	}
	c.Start()

	wg.Wait()
}

type TickerHandled struct {
	Time  int64
	Price *big.Float
}

func handleTickers(mu *sync.Mutex, wg *sync.WaitGroup, source internal.HttpStreamSubscriber, exchgRates map[string][]TickerHandled) {
	chTicker, chErr := source.SubscribePriceStream(internal.BTCUSDTicker)
LOOP:
	for {
		select {
		case ticker := <-chTicker:
			price, _, err := big.ParseFloat(ticker.Price, 10, 100, big.ToNearestAway)
			if err != nil {
				log.Printf("error parse float [%s]", err)
				continue LOOP
			}
			th := TickerHandled{
				Time:  ticker.Time.Unix(),
				Price: price,
			}

			mu.Lock()
			exchgRates[source.Url] = append(exchgRates[source.Url], th)
			mu.Unlock()
		case <-chErr:
			wg.Done()
			break LOOP
		}
	}
}

func handleJob(mu *sync.Mutex, exchgRates map[string][]TickerHandled) {
	var counter int64
	sum := new(big.Float)
	curr := time.Now().Unix()

	mu.Lock()
	for url, rates := range exchgRates {
		for _, rate := range rates {
			if curr-60 > rate.Time { // old rate
				continue
			}

			sum.Add(sum, rate.Price)
			counter++
		}
		delete(exchgRates, url)
	}
	mu.Unlock()

	if counter == 0 {
		return
	}

	fmt.Printf("%d, %s\n", curr, sum.Quo(sum, new(big.Float).SetInt64(counter)).Text('g', 15))
}
