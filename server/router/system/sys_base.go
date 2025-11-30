package system

import (
	"github.com/gin-gonic/gin"
	shopApi "github.com/flipped-aurora/gin-vue-admin/server/api/v1/shop"
)

type BaseRouter struct{}

func (s *BaseRouter) InitBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("base")
	productApi := shopApi.ProductApi{}
	{
		baseRouter.POST("buy", productApi.Buy)
		baseRouter.GET("find", productApi.Find)
		baseRouter.POST("login", baseApi.Login)
		baseRouter.POST("captcha", baseApi.Captcha)
	}
	return baseRouter
}
