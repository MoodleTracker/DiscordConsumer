package main

import (
	"context"
	"fmt"
	"github.com/MoodleTracker/Protocol-Go/protocol"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"log"
)

func main() {
	err := godotenv.Load()
	// epic go error handling
	if err != nil && err.Error() != "open .env: no such file or directory" {
		log.Fatal("Error loading .env file: ", err)
		return
	}

	// connect to kafka
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{GetEnv("KAFKA_BROKER")},
		Topic:     "events-upcoming",
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})

	// read messages
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("failed to read message:", err)
			break
		}

		// eww this sucks lmao
		// also I have literally no clue if this is the right way to do this or if this is future-proof but eh
		value := m.Value[6:]
		message := &protocol.UpcomingEvent{}
		err = proto.Unmarshal(value, message)
		if err != nil {
			log.Println("failed to unmarshal message:", err)
			continue
		}

		fmt.Printf("message at offset %d: %s = %#v\n", m.Offset, string(m.Key), message.String())
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
