package initialize

import (
	"context"
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	// "github.com/flipped-aurora/gin-vue-admin/server/mcp/client"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Mongo(){
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil{
		fmt.Println("failed to connect MongoDB:", err)
		return
	}
	err = client.Ping(ctx,nil)
	if err != nil{
		fmt.Println(err)
		return
	}
	global.GVA_MONGO = client
	fmt.Println("success to connect to MongoDB")
}