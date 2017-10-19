package main

import (
	"context"
	"fmt"
	"log"

	kafka "github.com/segmentio/kafka-go"
)

func main() {
	// to consume messages
	topic := "my-topic-3"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "127.0.0.1:9192", topic, partition)
	if err != nil {
		log.Fatalln(err)
	}

	offset1st, offsetlast, err := conn.ReadOffsets()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("first-offset: %v\n", offset1st)
	fmt.Printf("last-offset: %v\n", offsetlast)

	s, err := conn.Seek(10, 2)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("seek: %v\n", s)

	m, err := conn.ReadMessage(10000)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(m)

	batch := conn.ReadBatch(0, 1e6) // fetch 10KB min, 1MB max
	b := make([]byte, 10e3)         // 10KB max per message
	for {
		i, err := batch.Read(b)
		fmt.Println(i, err)
		if err != nil {
			break
		}
		fmt.Println(string(b))
	}

	batch.Close()
	conn.Close()
}
