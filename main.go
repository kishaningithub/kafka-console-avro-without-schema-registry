package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/uuid"
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
	TlsMode          string
	CertLocation     string
	KeyLocation      string
	CACertLocation   string
}

const (
	TLS_MODE_NONE = "NONE"
	TLS_MODE_TLS  = "TLS"
	TLS_MODE_MTLS = "MTLS"
)

func main() {
	appConfig := loadAppConfig()
	kafkaReaderConfig := kafka.ReaderConfig{
		Brokers:     appConfig.BootstrapServers,
		GroupID:     uuid.New().String(),
		Topic:       appConfig.Topic,
		StartOffset: kafka.LastOffset,
	}
	if appConfig.TlsMode != TLS_MODE_NONE {
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
		if len(message.Headers) > 0 {
			fmt.Println("Headers")
			fmt.Println("=====")
			for _, header := range message.Headers {
				fmt.Println(header.Key, string(header.Value))
			}
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
	tlsMode := flag.String("tls-mode", "NONE", "Valid values are NONE,TLS,MTLS")
	certLocation := flag.String("tls-cert", "", "certificate file location. Eg. /certs/cert.pem. Required if tls-mode is MTLS")
	keyLocation := flag.String("tls-key", "", "key file location. Eg. /certs/key.pem. Required if tls-mode is MTLS")
	caCertLocation := flag.String("tls-ca-cert", "", "CA cert file location. Eg. /certs/ca.pem. If not given system CA would take effect")
	flag.Parse()
	appConfig := Config{
		Topic:            *topic,
		BootstrapServers: strings.Split(*bootstrapServersCSV, ","),
		TlsMode:          *tlsMode,
		CertLocation:     *certLocation,
		KeyLocation:      *keyLocation,
		CACertLocation:   *caCertLocation,
	}
	_, _ = os.Stderr.WriteString(fmt.Sprintf("loaded config %+v \n", appConfig))
	return appConfig
}

func getTLSConfig(config Config) *tls.Config {
	tlsConfig := &tls.Config{
		RootCAs:    getCACert(config),
		MinVersion: tls.VersionTLS12,
	}
	if config.TlsMode == TLS_MODE_MTLS {
		certLocation := config.CertLocation
		keyLocation := config.KeyLocation
		cert, err := tls.LoadX509KeyPair(certLocation, keyLocation)
		if err != nil {
			panic(fmt.Errorf("error while loading cert %s and key %s: %w", certLocation, keyLocation, err))
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	return tlsConfig
}

func getCACert(config Config) *x509.CertPool {
	caCertLocation := config.CACertLocation
	if caCertLocation == "" {
		return nil
	}
	caCert, err := ioutil.ReadFile(caCertLocation)
	if err != nil {
		panic(fmt.Errorf("error while loading ca cert from %s: %w", caCertLocation, err))
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return caCertPool
}
