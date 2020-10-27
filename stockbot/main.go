package stockbot

import (
	"encoding/json"
	"fmt"
	"time"

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
			fmt.Println("message", req)
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

func getStock(cmd string) (string, error) {
	t := time.Now()
	return "APPL.US quote is $93.42 per share." + t.Format("20060102150405"), nil
}
