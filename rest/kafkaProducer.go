package rest

import (
	"github.com/Shopify/sarama"
	"github.com/commitHub/commitBlockchain/wire"
)

// NewProducer is a producer to send messages to kafka
func NewProducer(kafkaPorts []string) sarama.SyncProducer {
	producer, err := sarama.NewSyncProducer(kafkaPorts, nil)
	if err != nil {
		panic(err)
	}
	return producer
}

//KafkaProducerDeliverMessage : delivers messages to kafka
func KafkaProducerDeliverMessage(msg KafkaMsg, topic string, producer sarama.SyncProducer, cdc *wire.Codec) error {

	kafkaStoreBytes, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	sendmsg := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(kafkaStoreBytes),
	}
	_, _, err = producer.SendMessage(&sendmsg)
	if err != nil {
		return err
	}
	return nil
}

//SendToKafka : handles sending message to kafka
func SendToKafka(msg KafkaMsg, kafkaState KafkaState, cdc *wire.Codec) []byte {
	err := KafkaProducerDeliverMessage(msg, "Topic", kafkaState.Producer, cdc)
	if err != nil {
		jsonResponse, err := cdc.MarshalJSON(struct {
			Response string `json:"response"`
		}{Response: "Something is up with kafka server, restart rest and kafka."})
		if err != nil {
			panic(err)
		}
		SetTicketIDtoDB(msg.TicketID, kafkaState.KafkaDB, cdc, jsonResponse)
	} else {
		jsonResponse, err := cdc.MarshalJSON(struct {
			Response string `json:"response"`
		}{Response: "Request in process, wait and try after some time"})
		if err != nil {
			panic(err)
		}
		SetTicketIDtoDB(msg.TicketID, kafkaState.KafkaDB, cdc, jsonResponse)
	}
	jsonResponse, err := cdc.MarshalJSON(struct {
		TicketID Ticket `json:"ticketID"`
	}{TicketID: msg.TicketID})
	if err != nil {
		panic(err)
	}
	return jsonResponse
}
