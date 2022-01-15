package internal

import (
	"fmt"
	"math/rand"
	"time"
)

type HttpStreamSubscriber struct {
	Url string
}

func (s *HttpStreamSubscriber) SubscribePriceStream(ticker Ticker) (chan TickerPrice, chan error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	chErr := make(chan error)
	chTicker := make(chan TickerPrice)

	go func() {
		// REAL
		/*resp, err := http.Get(s.Url)
		if err != nil {
			chErr <- fmt.Errorf("connecting error has occurred")
			close(chTicker)
			close(chErr)
			return
		}

		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				chErr <- fmt.Errorf("getting data error has occurred")
				close(chTicker)
				close(chErr)
				break
			}

			tp := TickerPrice{}
			err = json.Unmarshal(line, &tp)
			if err != nil {
				chErr <- fmt.Errorf("unmashalling data error has occurred")
				close(chTicker)
				close(chErr)
				break
			}

			chTicker <- tp
		}*/

		// FAKE
		priceTicker := time.NewTicker(5 * time.Second)
		errTicker := time.NewTicker(130 * time.Second)

	LOOP:
		for {
			select {
			case <-priceTicker.C:
				tp := TickerPrice{
					Ticker: BTCUSDTicker,
					Time:   time.Now(),
					Price:  fmt.Sprintf("%f", r.Float64()),
				}
				chTicker <- tp
			case <-errTicker.C:
				chErr <- fmt.Errorf("getting data error has occurred")

				close(chTicker)
				close(chErr)

				break LOOP
			}
		}
	}()

	return chTicker, chErr
}
