package stockbot

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	csvtag "github.com/artonge/go-csv-tag/v2"
	"github.com/carlware/gochat/common/config"
	"github.com/carlware/gochat/common/mq/interfaces/rabbitmq"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Request struct {
	Command string `json:"command"`
	Extra   string `json:"extra"`
}

type Response struct {
	Result string `json:"result"`
	Extra  string `json:"extra"`
	Error  *Error `json:"error"`
}

func Run(cfg *config.Configuration) {
	mq, err := rabbitmq.NewServer(cfg.RabbiMQ.Host)
	if err != nil {
		panic(err)
	}

	// listen
	go func() {
		msgs, _ := mq.Listen("request")
		for msg := range msgs {
			req := Request{}
			json.Unmarshal(msg, &req)
			result, err := getStock(req.Command)
			response := Response{}
			response.Result = result
			response.Extra = req.Extra
			if err != nil {
				response.Error = &Error{
					Code:    "command",
					Message: err.Error(),
				}
			}
			encoded, _ := json.Marshal(response)
			mq.Send(encoded, "results")
		}
	}()

	forever := make(chan bool)
	<-forever
}

type Stock struct {
	Symbol string  `csv:"Symbol"`
	Date   string  `csv:"Date"`
	Time   string  `csv:"Time"`
	Open   float64 `csv:"Open"`
	High   float64 `csv:"High"`
	Low    float64 `csv:"Low"`
	Close  float64 `csv:"Close"`
	Volume string  `csv:"Volume"`
}

func getStock(cmd string) (string, error) {
	stockCode := strings.ToLower(cmd)
	query := fmt.Sprintf("https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv", stockCode)

	resp, err := http.Get(query)
	if err != nil {
		return "", errors.New("server error")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("Error from Stooq server")
	}

	stocks := []Stock{}
	err = csvtag.LoadFromReader(resp.Body, &stocks)

	if len(stocks) < 1 {
		return "", errors.New(fmt.Sprintf("Stock: %s not available", cmd))
	}

	return fmt.Sprintf("%s quote is $%.2f per share.", cmd, stocks[0].Close), nil
}
