package shop

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
)

type OrderMessage struct{
	ProductID uint `json:"product_id"`
	Count int `json:"count"`
}

func StartOrderConsumer(){
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"},nil)
	if err != nil{
		fmt.Println("consumer start failed:", err)
		return
	}
	partitionList, err := consumer.Partitions("orders")
	if err != nil{
		fmt.Println("unable to get the partition:", err)
		return
	}
	fmt.Println("consumer started,listening")
	for _, partition := range partitionList{
		pc,_ := consumer.ConsumePartition("orders", partition, sarama.OffsetNewest)
		go func(pc sarama.PartitionConsumer){
			defer pc.Close()
			var productService ProductService
			for message := range pc.Messages(){
				var msg OrderMessage
				json.Unmarshal(message.Value, &msg)
				fmt.Printf("order received: merchandiseID = %d, numbers = %d\n", msg.ProductID, msg.Count)
				err := productService.ReduceStockOptimistic(msg.ProductID, msg.Count)
				if err != nil{
					fmt.Println(err)
				}else{
					fmt.Println("done!")
				}
			}
		}(pc)
	}
}