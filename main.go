package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/linkedin/goavro"
	"github.com/segmentio/kafka-go"
	"strings"
)

func main() {
	topic := flag.String("topic", "topic name", "The topic id to consume on.")
	bootstrapServersCSV := flag.String("bootstrap-servers", "server(s) to connect to", "REQUIRED: The server(s) to connect to.")

	flag.Parse()

	bootstrapServers := strings.Split(*bootstrapServersCSV, ",")
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: bootstrapServers,
		GroupID: "GroupID",
		Topic:   *topic,
	})
	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			panic(fmt.Errorf("error while reading message from kafka topic %s: %w", *topic, err))
		}
		ocfReader, err := goavro.NewOCFReader(bytes.NewBuffer(message.Value))
		if err != nil {
			panic(fmt.Errorf("data is not in plain avro format %s: %w", *topic, err))
		}
		avroSchema := ocfReader.Codec().Schema()
		fmt.Println("Schema")
		fmt.Println("=====")
		var parsedSchema map[string]interface{}
		err = json.Unmarshal([]byte(avroSchema), &parsedSchema)
		if err != nil {
			panic(fmt.Errorf("error while unmarshalling avro schema %s: %w", avroSchema, err))
		}
		jsonSchemaBytes, err := json.MarshalIndent(parsedSchema, "", "  ")
		if err != nil {
			panic(fmt.Errorf("error while marshalling avro schema %s: %w", avroSchema, err))
		}
		fmt.Println(string(jsonSchemaBytes))
		fmt.Println("Data")
		fmt.Println("=====")
		for ocfReader.Scan() {
			record, _ := ocfReader.Read()
			jsonRecord, err := json.MarshalIndent(record, "", "  ")
			if err != nil {
				panic(fmt.Errorf("unable to marshal record %v as json: %w", record, err))
			}
			fmt.Println(string(jsonRecord))
		}
		fmt.Println()
	}
}
