package main

import (
	"encoding/json"
	"fmt"
	"sync"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/netoudi/codeinvest/stock-exchange/internal/infra/kafka"
	"github.com/netoudi/codeinvest/stock-exchange/internal/market/dto"
	"github.com/netoudi/codeinvest/stock-exchange/internal/market/entity"
	"github.com/netoudi/codeinvest/stock-exchange/internal/market/transformer"
)

func main() {
	ordersIn := make(chan *entity.Order)
	ordersOut := make(chan *entity.Order)
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	kafkaMsgChan := make(chan *ckafka.Message)
	configMapConsumer := &ckafka.ConfigMap{
		"bootstrap.servers": "host.docker.internal:9094",
		"group.id":          "myGroup",
		"auto.offset.reset": "latest",
	}
	configMapProducer := &ckafka.ConfigMap{
		"bootstrap.servers": "host.docker.internal:9094",
	}
	producer := kafka.NewProducer(configMapProducer)
	consumer := kafka.NewConsumer(configMapConsumer, []string{"input"})

	go consumer.Consume(kafkaMsgChan)

	book := entity.NewBook(ordersIn, ordersOut, wg)
	go book.Trade()

	go func() {
		for msg := range kafkaMsgChan {
			fmt.Println(msg)
			wg.Add(1)
			fmt.Println(string(msg.Value))
			tradeInput := dto.TradeInput{}
			err := json.Unmarshal(msg.Value, &tradeInput)
			if err != nil {
				panic(err)
			}
			order := transformer.TransformInput(tradeInput)
			ordersIn <- order
		}
	}()

	fmt.Println("🚀 App is running...")
	for res := range ordersOut {
		output := transformer.TransformOutput(res)
		outputJson, err := json.MarshalIndent(output, "", " ")
		fmt.Println(string(outputJson))
		if err != nil {
			fmt.Println(err)
		}
		err = producer.Publish(outputJson, []byte("orders"), "output")
		if err != nil {
			fmt.Println(err)
		}
	}
}
