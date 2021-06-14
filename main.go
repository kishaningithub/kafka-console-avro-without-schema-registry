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
	"os"
	"strings"
	"time"
)

type Config struct {
	Topic            string
	BootstrapServers []string
	TlsEnabled       bool
	CertLocation     string
	KeyLocation      string
	CACertLocation   string
}

func main() {
	appConfig := loadAppConfig()
	kafkaReaderConfig := kafka.ReaderConfig{
		Brokers: appConfig.BootstrapServers,
		GroupID: "GroupID",
		Topic:   appConfig.Topic,
	}
	if appConfig.TlsEnabled {
		tlsConfig := getTLSConfig(appConfig)
		kafkaReaderConfig.Dialer = &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
			TLS:       tlsConfig,
		}
	}
	reader := kafka.NewReader(kafkaReaderConfig)
	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			panic(fmt.Errorf("error while reading message from kafka topic %s: %w", appConfig.Topic, err))
		}
		ocfReader, err := goavro.NewOCFReader(bytes.NewBuffer(message.Value))
		if err != nil {
			panic(fmt.Errorf("data is not in plain avro format %s: %w", appConfig.Topic, err))
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

func loadAppConfig() Config {
	topic := flag.String("topic", "", "REQUIRED: The topic id to consume on.")
	bootstrapServersCSV := flag.String("bootstrap-servers", "", "REQUIRED: The server(s) to connect to.")
	tlsEnabled := flag.Bool("tls-enabled", false, "If this is set to true, then the other tls flags are required")
	certLocation := flag.String("tls-cert", "", "certificate file location. Eg. /certs/cert.pem")
	keyLocation := flag.String("tls-key", "", "key file location. Eg. /certs/key.pem")
	caCertLocation := flag.String("tls-ca-cert", "", "CA cert file location. Eg. /certs/ca.pem")
	flag.Parse()
	appConfig := Config{
		Topic:            *topic,
		BootstrapServers: strings.Split(*bootstrapServersCSV, ","),
		TlsEnabled:       *tlsEnabled,
		CertLocation:     *certLocation,
		KeyLocation:      *keyLocation,
		CACertLocation:   *caCertLocation,
	}
	_, _ = os.Stderr.WriteString(fmt.Sprintf("loaded config %+v \n", appConfig))
	return appConfig
}

func getTLSConfig(config Config) *tls.Config {
	certLocation := config.CertLocation
	keyLocation := config.KeyLocation
	caCertLocation := config.CACertLocation
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
