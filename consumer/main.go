package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

func ConnectBroker(brokerURL []string) (sarama.Consumer, error) {

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	return sarama.NewConsumer(brokerURL, config)

}
func main() {

	topic := "coffee_order"
	msgCnt := 0
	// create a new consumer
	worker, err := ConnectBroker([]string{"localhost:9092"})
	if err != nil {
		panic(err)
	}

	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	// handle os signals
	sigchan := make(chan os.Signal, 1)
	doneCh := make(chan struct{})

	signal.Notify(sigchan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				msgCnt++
				fmt.Println(msgCnt)
				order := string(msg.Value)
				fmt.Println("Brewing coffee", order)
			case <-sigchan:
				fmt.Println("\nInterrupt signal detected in worker...")
				close(doneCh)
				return
			}
		}
	}()
	fmt.Println("Your server is running. Press Ctrl+C to exit.")

	// Wait for the goroutine to signal it's done
	<-doneCh

	// close the consumer

	fmt.Println("Processed", msgCnt, "messages")

	if err := worker.Close(); err != nil {
		panic(err)
	}

}
