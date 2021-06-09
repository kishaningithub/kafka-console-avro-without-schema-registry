# kafka-console-avro-without-schema-registry

Tail kafka avro topic data without confluent schema registry overhead

## Installation

```shell
brew tap kishaningithub/tap
brew install kafka-console-avro-without-schema-registry
```

## Usage

```shell
kafka-console-avro-without-schema-registry --topic stats --bootstrap-servers localhost:9092
```