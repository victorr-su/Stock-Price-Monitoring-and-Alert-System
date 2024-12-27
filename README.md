# Commands to start kafka locally

`bin/zookeeper-server-start.sh config/zookeeper.properties`

`bin/kafka-server-start.sh config/server.properties`

# Start program locally

`go run cmd/app/main.go`

# Manually check for kafka messages

`bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic stock-prices --from-beginning`
