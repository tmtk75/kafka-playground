package main

import (
	"context"
	"fmt"
	"log"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func main() {
	topic := "my-topic-3"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "127.0.0.1:9192", topic, partition)
	if err != nil {
		log.Fatalln(err)
	}
	offset1st, offsetlast, err := conn.ReadOffsets()
	fmt.Printf("first-offset: %v, last-offset: %v\n", offset1st, offsetlast)

	// mhhh... I cannot understand how it works.
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"127.0.0.1:9192", "127.0.0.1:9292", "127.0.0.1:9392"},
		Topic:     topic,
		Partition: partition,
		MinBytes:  10,   // 10KB
		MaxBytes:  10e6, // 10MB
	})
	r.SetOffset(offsetlast - 10)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		m, err := r.ReadMessage(ctx)
		if err != nil {
			log.Println(err)
			break
		}
		//fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		fmt.Printf("%v\n", m.Offset)
	}

	r.Close()
}
