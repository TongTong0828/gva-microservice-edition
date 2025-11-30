package initialize

import (
	"fmt"

	"github.com/IBM/sarama"
	// "github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	// "github.com/flipped-aurora/gin-vue-admin/server/mcp/client"
)

func Kafka(){
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"},config)
	if err != nil{
		fmt.Println("kafka failed:", err)
		return
	}
	global.GVA_KAFKA_PRODUCER = client
	fmt.Println("kafka success!")
}