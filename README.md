# gokafka

This is simple example of Golang and Kafka Integration

## How to run

1.Run the docker compose 
```docker
docker compose up -d
```

2.Create the topic on Kafka 
```docker
docker exec -it kafka-like /bin/bash
```
Then create a topic (example : user-topic)
```bash
kafka-topics.sh --create --topic user-topic --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
```
Check the topic has been created?
```
kafka-topics.sh --list --bootstrap-server localhost:9092
```
3.Run `go mod tidy` to get the libraries that needed

4.Run the customer code :
```
cd customer
go run customer.go
```

5.Run the producer code :
```
cd producer
go run producer.go
```

## Concurrency Pattern
The producer.go using Pipeline Concurrency Pattern to handle the message. We can see the code : 
```go
chnListUser: = make(chan Message)
// Get the message then send into `chnListUser`
go func() {
    defer close(chnListUser)
    for i: = 0;
    i < len(userId);
    i++{
        message: = Message {
            UserId: userId[i],
            Name: name[i],
            Score: score[i],
        }

            chnListUser < -message
    }
}()

chnFilterPassedScore: = make(chan[] byte, len(userId))
// Proceed the `chnListUser` then filter the Score > 50
// Assign to channel `chnFilterPassedScore`
go func() {
    defer close(chnFilterPassedScore)
    for val: = range chnListUser {
        if val.Score > 50 {
            jsonMessage, err: = json.Marshal(val)
            chnFilterPassedScore < -jsonMessage

            if err != nil {
                log.Fatalln(err.Error())
                os.Exit(1)
            }
        }
    }
}()

wg: = sync.WaitGroup {}

// Get the latest/final value that has been filtered
// Send to Kafka producer
wg.Add(1)
go func() {
    defer wg.Done()
    for k: = range chnFilterPassedScore {
        msg: = & sarama.ProducerMessage {
            Topic: conf.TopicName,
            Value: sarama.StringEncoder(k),
        }

            log.Println("Sending message :", msg)
        _,
        _,
        err = producer.SendMessage(msg)
        if err != nil {
            log.Fatalln(err.Error())
            os.Exit(1)
        }

        log.Println("Message sent!")
    }
}()
```
