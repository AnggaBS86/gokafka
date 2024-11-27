package main

import (
	"context"
	"fmt"
	"gokafka/conf"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type consumerGroupHandler struct{}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message received : %s\n", string(msg.Value))
		fmt.Println()
		sess.MarkMessage(msg, "")
	}
	return nil
}

func main() {
	brokers := []string{conf.Host}
	groupID := conf.GroupId

	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Panicf("Error occured when creating consumer group client: %s", err.Error())
		fmt.Println()
	}

	ctx := context.Background()
	for {
		err := consumerGroup.Consume(ctx, []string{conf.TopicName}, &consumerGroupHandler{})
		if err != nil {
			log.Fatalf("Error from consumer: %s", err.Error())
			fmt.Println()
		}
	}
}
