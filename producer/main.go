package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	conf "gokafka/conf"

	"github.com/IBM/sarama"
)

type Message struct {
	UserId int    `json:"user_id"`
	Name   string `json:"name"`
	Score  int    `json:"score"`
}

func main() {
	brokers := []string{conf.Host}
	producer, err := sarama.NewSyncProducer(brokers, nil)
	if err != nil {
		log.Fatalln("Error occured when starting Sarama producer:", err)
		os.Exit(1)
	}

	// This is fake/dummy data
	userId := [5]int{100001, 100002, 100003, 100004, 100005}
	name := [5]string{"John", "Doe", "Anderson", "Andre", "Michelle"}
	score := [5]int{90, 80, 50, 70, 40}

	chnListUser := make(chan Message)
	// Get the message then send into `chnListUser`
	go func() {
		defer close(chnListUser)
		for i := 0; i < len(userId); i++ {
			message := Message{
				UserId: userId[i],
				Name:   name[i],
				Score:  score[i],
			}

			chnListUser <- message
		}
	}()

	chnFilterPassedScore := make(chan []byte, len(userId))
	// Proceed the `chnListUser` then filter the Score > 50
	// Assign to channel `chnFilterPassedScore`
	go func() {
		defer close(chnFilterPassedScore)
		for val := range chnListUser {
			if val.Score > 50 {
				jsonMessage, err := json.Marshal(val)
				chnFilterPassedScore <- jsonMessage

				if err != nil {
					log.Fatalln(err.Error())
					os.Exit(1)
				}
			}
		}
	}()

	wg := sync.WaitGroup{}

	// Get the latest/final value that has been filtered
	// Send to Kafka producer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for k := range chnFilterPassedScore {
			msg := &sarama.ProducerMessage{
				Topic: conf.TopicName,
				Value: sarama.StringEncoder(k),
			}

			log.Println("Sending message :", msg)
			_, _, err = producer.SendMessage(msg)
			if err != nil {
				log.Fatalln(err.Error())
				os.Exit(1)
			}

			log.Println("Message sent!")
		}
	}()

	wg.Wait()
}
