package utils

import (
	"context"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type MongoWriter struct{}

func (mw *MongoWriter) Write(p []byte) (n int, err error){
	if global.GVA_MONGO == nil{
		return len(p), nil
	}
	type LogEntry struct{
		Content string `bson:"content"`
		CreatedAt time.Time `bson:"creater_at"`
	}
	go func(logMsg string){
		collection := global.GVA_MONGO.Database("gav_logs").Collection("system_logs")
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		collection.InsertOne(ctx, LogEntry{
			Content: logMsg,
			CreatedAt: time.Now(),
		})
	}(string(p))
	return len(p), nil
}

func (mw *MongoWriter) Sync() error{
	return nil
}