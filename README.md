# Kafka console avro without schema registry
![build workflow](https://github.com/kishaningithub/kafka-console-avro-without-schema-registry/actions/workflows/build.yml/badge.svg)

Tail kafka avro topic data without confluent schema registry overhead

This expects the data to be written in [Object Container File (OCF)](https://avro.apache.org/docs/current/spec.html#Object+Container+Files) format

## Installation

```shell
$ brew install kishaningithub/tap/kafka-console-avro-without-schema-registry
```

## Upgrading

```shell
$ brew upgrade kafka-console-avro-without-schema-registry
```

## Usage

```shell
$ kafka-console-avro-without-schema-registry -help
Usage of kafka-console-avro-without-schema-registry:
  -bootstrap-servers string
    	REQUIRED: The server(s) to connect to.
  -tls-ca-cert string
    	CA cert file location. Eg. /certs/ca.pem. If not given system CA would take effect
  -tls-cert string
    	certificate file location. Eg. /certs/cert.pem. Required if tls-mode is MTLS
  -tls-key string
    	key file location. Eg. /certs/key.pem. Required if tls-mode is MTLS
  -tls-mode string
    	Valid values are NONE,TLS,MTLS (default "NONE")
  -topic string
    	REQUIRED: The topic id to consume on.
```

## Examples

### Basic usage

```shell
$ kafka-console-avro-without-schema-registry --topic example --bootstrap-servers localhost:9092
Headers
=====
machineid 697d86e5-feea-35f7-a988-123f7614e4be
Schema
=====
{"fields":[{"name":"time","type":"long"},{"default":"","description":"Process id","name":"process_id","type":"string"}],"name":"example","namespace":"com.example","type":"record","version":1}
Data
=====
{"time":1617104831727, "process_id":"ID1"}
{"time":1717104831727, "process_id":"ID2"}

Schema
=====
{"fields":[{"name":"time","type":"long"}],"name":"example","namespace":"com.example","type":"record","version":2}
Data
=====
{"time":1817104831727}
{"time":1917104831727}
```

### With TLS

```shell
$ kafka-console-avro-without-schema-registry --topic example --bootstrap-servers localhost:9092 -tls-mode TLS -tls-ca-cert /certs/ca.pem
Headers
=====
machineid 697d86e5-feea-35f7-a988-123f7614e4be
Schema
=====
{"fields":[{"name":"time","type":"long"},{"default":"","description":"Process id","name":"process_id","type":"string"}],"name":"example","namespace":"com.example","type":"record","version":1}
Data
=====
{"time":1617104831727, "process_id":"ID1"}
{"time":1717104831727, "process_id":"ID2"}

Schema
=====
{"fields":[{"name":"time","type":"long"}],"name":"example","namespace":"com.example","type":"record","version":2}
Data
=====
{"time":1817104831727}
{"time":1917104831727}
```

### With MTLS

```shell
$ kafka-console-avro-without-schema-registry --topic example --bootstrap-servers localhost:9092 -tls-mode MTLS -tls-cert /certs/cert.pem -tls-key /certs/key.pem -tls-ca-cert /certs/ca.pem
Headers
=====
machineid 697d86e5-feea-35f7-a988-123f7614e4be
Schema
=====
{"fields":[{"name":"time","type":"long"},{"default":"","description":"Process id","name":"process_id","type":"string"}],"name":"example","namespace":"com.example","type":"record","version":1}
Data
=====
{"time":1617104831727, "process_id":"ID1"}
{"time":1717104831727, "process_id":"ID2"}

Schema
=====
{"fields":[{"name":"time","type":"long"}],"name":"example","namespace":"com.example","type":"record","version":2}
Data
=====
{"time":1817104831727}
{"time":1917104831727}
```
