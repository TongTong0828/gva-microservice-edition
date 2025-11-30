package shop

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	// "github.com/qiniu/go-sdk/v7/internal/cache"

	// "github.com/shoenig/test/interfaces"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	// "github.com/flipped-aurora/gin-vue-admin/server/global"
	// "github.com/flipped-aurora/gin-vue-admin/server/model/shop"
)

type ProductService struct{}

func (s *ProductService) ReduceStockOptimistic(productID uint, buyCount int)(err error){
	var product shop.Product
	if err := global.GVA_DB.First(&product, productID).Error; err != nil{
		return errors.New("not existing product")
	}
	if product.Stock < buyCount{
		return errors.New("not enough stock")
	}
	result := global.GVA_DB.Model(&shop.Product{}).
	Where("id = ? AND version = ? AND stock >= ?", productID, product.Version, buyCount).
	Updates(map[string]interface{}{
		"stock": gorm.Expr("stock - ?", buyCount),
		"version": gorm.Expr("version + 1"),
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("purchase failed, data changed or not enough stock")
	}
	return nil
}

func (s *ProductService) GetProductCached(id uint) (shop.Product, error){
	var product shop.Product
	cacheKey := fmt.Sprintf("product:%d", id)
	val, err := global.GVA_REDIS.Get(context.Background(),cacheKey).Result()
	if err == nil{
		json.Unmarshal([]byte(val), &product)
		fmt.Println("redis hit")
		return product, nil
	}
	fmt.Println("loading mysql")
	if err := global.GVA_DB.First(&product, id).Error; err != nil{
		return product, err
	}
	jsonBytes, _ := json.Marshal(product)
	global.GVA_REDIS.Set(context.Background(),cacheKey,jsonBytes,10 * time.Minute)
	return product, nil
}