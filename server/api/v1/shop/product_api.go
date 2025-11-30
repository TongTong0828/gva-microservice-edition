package shop

import (
	"encoding/json"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	shopService "github.com/flipped-aurora/gin-vue-admin/server/service/shop" // 别名
	"github.com/gin-gonic/gin"

	// "github.com/flipped-aurora/gin-vue-admin/server/global"
	// "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service/shop" // 引用 Message 结构体
	// "github.com/gin-gonic/gin"
)

type ProductApi struct{}

var ProductService = shopService.ProductService{}

func(p *ProductApi) Buy(c *gin.Context){
	productId := 1
	count := 1
	msg := shop.OrderMessage{
		ProductID: uint(productId),
		Count: count,
	}
	msgBytes, _ := json.Marshal(msg)
	kafkaMsg := &sarama.ProducerMessage{
		Topic: "orders",
		Value: sarama.ByteEncoder(msgBytes),
	}
	_, _, err := global.GVA_KAFKA_PRODUCER.SendMessage(kafkaMsg)
	if err != nil{
		global.GVA_LOG.Error("failed" + err.Error())
		response.FailWithMessage("busy, try later", c)
		return
	}
	response.OkWithMessage("received, proceeding", c)
}

func (p *ProductApi) Find(c *gin.Context){
	idStr := c.Query("id")
	ProductId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil{
		return
	}
	product, err := ProductService.GetProductCached(uint(ProductId))
	if err != nil{
		response.FailWithMessage("query failed", c)
		return
	}
	response.OkWithData(product, c)
}