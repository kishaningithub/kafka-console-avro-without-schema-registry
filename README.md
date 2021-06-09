# Kafka console avro without schema registry

Tail kafka avro topic data without confluent schema registry overhead

This expects the data to be written in [Object Container File (OCF)](https://avro.apache.org/docs/current/spec.html#Object+Container+Files) format

## Installation

```shell
brew tap kishaningithub/tap
brew install kafka-console-avro-without-schema-registry
```

## Upgrading

```shell
brew upgrade kafka-console-avro-without-schema-registry
```

## Usage

```shell
kafka-console-avro-without-schema-registry --topic example --bootstrap-servers localhost:9092
```

## Output

```shell
$ kafka-console-avro-without-schema-registry --topic example --bootstrap-servers localhost:9092

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