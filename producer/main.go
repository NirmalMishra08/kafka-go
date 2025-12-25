package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/IBM/sarama"
)

type Order struct {
	ConsumerName string `json:"consumer_name,omitempty"`
	CoffeeType   string `json:"coffee_type,omitempty"`
}

func main() {
	http.HandleFunc("/order", placeOrder)

	http.ListenAndServe(":8000", nil)
}

func ConnectBroker(brokerURL []string) (sarama.SyncProducer, error) {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 4

	return sarama.NewSyncProducer(brokerURL, config)

}
func sendPushOrderToQueue(topic string, message []byte) error {
	brokerURL := []string{"localhost:9092"}

	// connect connection
	producer, err := ConnectBroker(brokerURL)
	if err != nil {
		return err
	}

	defer producer.Close()

	// create a new kafka message

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Order of topic : %s \n partition : %s \n offset : %s", topic, partition, offset)

	return nil

}

func placeOrder(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "invalid request type", http.StatusBadRequest)
	}

	var order Order

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "not able to parse body", http.StatusBadRequest)
		return
	}

	orderInBytes, err := json.Marshal(order)
	if err != nil {
		http.Error(w, "not able to parse string to bytes", http.StatusBadRequest)
		return
	}

	// send bytes
	err = sendPushOrderToQueue("coffee_order", orderInBytes)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	// respond back to the user;

	response := map[string]interface{}{
		"successs": true,
		"msg":      "Order for Customer:" + order.ConsumerName + "send successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

}
