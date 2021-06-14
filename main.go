package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/linkedin/goavro"
	"github.com/segmentio/kafka-go"
	"io/ioutil"
	"strings"
	"time"
)

func main() {
	topic := flag.String("topic", "", "REQUIRED: The topic id to consume on.")
	bootstrapServersCSV := flag.String("bootstrap-servers", "", "REQUIRED: The server(s) to connect to.")
	tlsEnabled := flag.Bool("tls-enabled", false, "If this is set to true, then the other tls flags are required")
	certLocation := flag.String("tls-cert", "", "certificate file location. Eg. /certs/cert.pem")
	keyLocation := flag.String("tls-key", "", "key file location. Eg. /certs/key.pem")
	caCertLocation := flag.String("tls-ca-cert", "", "CA cert file location. Eg. /certs/ca.pem")
	flag.Parse()
	bootstrapServers := strings.Split(*bootstrapServersCSV, ",")
	config := kafka.ReaderConfig{
		Brokers: bootstrapServers,
		GroupID: "GroupID",
		Topic:   *topic,
	}
	if *tlsEnabled {
		tlsConfig := getTLSConfig(*certLocation, *keyLocation, *caCertLocation)
		config.Dialer = &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
			TLS:       tlsConfig,
		}
	}
	reader := kafka.NewReader(config)
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
		fmt.Println(avroSchema)
		fmt.Println("Data")
		fmt.Println("=====")
		for ocfReader.Scan() {
			record, _ := ocfReader.Read()
			jsonRecord, err := json.Marshal(record)
			if err != nil {
				panic(fmt.Errorf("unable to marshal record %v as json: %w", record, err))
			}
			fmt.Println(string(jsonRecord))
		}
		fmt.Println()
	}
}

func getTLSConfig(certLocation, keyLocation, caCertLocation string) *tls.Config {
	cert, err := tls.LoadX509KeyPair(certLocation, keyLocation)
	if err != nil {
		panic(fmt.Errorf("error while loading cert %s and key %s: %w", certLocation, keyLocation, err))
	}
	caCert, err := ioutil.ReadFile(caCertLocation)
	if err != nil {
		panic(fmt.Errorf("error while loading ca cert from %s: %w", caCertLocation, err))
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}
}
