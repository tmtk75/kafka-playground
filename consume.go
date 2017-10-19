package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Shopify/sarama"
)

func main() {
	Consume("", 0, 0)
}

func Consume(dpt string, offset int64, partition int) error {
	Consumer := func(m *Message) {
		fmt.Println(m)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Kill, os.Interrupt)
	go func() {
		<-signals
		cancel()
	}()

	return StartConsuming(ctx, Consumer)
}

type Message struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Topic     string    `json:"topic"`
	Partition int32     `json:"partition"`
	Offset    int64     `json:"offset"`
	Timestamp time.Time `json:"timestamp"` // only set if kafka is version 0.10+
}

func (m *Message) fill(c *sarama.ConsumerMessage) {
	m.Key = string(c.Key)
	m.Value = string(c.Value)
	m.Topic = c.Topic
	m.Partition = c.Partition
	m.Offset = c.Offset
	m.Timestamp = c.Timestamp
}

func StartConsuming(ctx context.Context, consumer func(m *Message)) error {
	sarama.Logger = log.New(os.Stdout, "[Sarama] ", log.LstdFlags)

	conf := sarama.NewConfig()
	conf.ClientID = "foobar"

	Hosts := []string{"127.0.0.1:9192", "127.0.0.1:9292", "127.0.0.1:9392"}
	Topic := "my-topic-3"
	Partition := 0
	Offset := sarama.OffsetNewest

	con, err := sarama.NewConsumer(Hosts, conf)
	if err != nil {
		return err
	}
	log.Printf("New consumer. hosts: %v", Hosts)

	pc, err := con.ConsumePartition(Topic, int32(Partition), int64(Offset))
	if err != nil {
		return err
	}
	log.Printf("Consume partition. topic: %v, partition: %v, offset: %v", Topic, Partition, Offset)

	go func() {
		for msg := range pc.Messages() {
			m := Message{}
			m.fill(msg)
			consumer(&m)
		}
	}()

	select {
	case <-ctx.Done():
		pc.AsyncClose()
	}

	if err := con.Close(); err != nil {
		fmt.Println("Failed to close consumer: ", err)
	}

	return nil
}
