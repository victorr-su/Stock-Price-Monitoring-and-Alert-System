package main

import (
	"Stock-Price-Monitoring-and-Alert-System/internal/kafka/consumer"
	"Stock-Price-Monitoring-and-Alert-System/internal/kafka/producer"
	"fmt"
)

func main(){
	fmt.Println("Stock price application begins");
	
	go producer.StartProducer()
	
	consumer.StartConsumer()
}