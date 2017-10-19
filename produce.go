package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Shopify/sarama"
)

func main() {
	var (
		max         = flag.Int("max", 1, "Max number to sent message")
		intervalStr = flag.String("interval", "1s", "Interval to send next message")
		topicName   = flag.String("topic", "my-topic-3", "Topic name")
	)
	flag.Parse()
	interval, err := time.ParseDuration(*intervalStr)
	if err != nil {
		log.Fatalln("%v", err)
	}

	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	conf.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	conf.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	conf.Producer.Return.Successes = true

	brokerList := []string{"127.0.0.1:9192", "127.0.0.1:9292", "127.0.0.1:9392"}

	producer, err := sarama.NewAsyncProducer(brokerList, conf)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Kill, os.Interrupt)
	go func() {
		<-signals
		cancel()
	}()

	go func() {
		for err := range producer.Errors() {
			log.Println(err)
		}
		cancel()
	}()

	go func() {
		for e := range producer.Successes() {
			log.Println(e.Value)
		}
		cancel()
	}()

	var start, end time.Time
	go func() {
		start = time.Now()
		for i := 0; i < *max && ctx.Err() == nil; i++ {
			producer.Input() <- &sarama.ProducerMessage{
				Topic: *topicName,
				//Key:   sarama.StringEncoder("a-key"),
				Value: &Msg{Val: fmt.Sprintf("Hello, World! %v", time.Now())},
			}
			time.Sleep(interval)
		}
		end = time.Now()
		producer.AsyncClose()
	}()

	select {
	case <-ctx.Done():
		log.Printf("done to send all: %v", time.Duration(end.UnixNano()-start.UnixNano()))
	}
}

type Msg struct {
	Val string
}

func (m *Msg) Encode() ([]byte, error) {
	return []byte(m.Val), nil
}

func (m *Msg) Length() int {
	return len([]byte(m.Val))
}
