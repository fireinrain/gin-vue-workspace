package cfscan

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AsnInfoRouter struct{}

// InitAsnInfoRouter 初始化 asnInfo表 路由信息
func (s *AsnInfoRouter) InitAsnInfoRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	asnInfoRouter := Router.Group("asnInfo").Use(middleware.OperationRecord())
	asnInfoRouterWithoutRecord := Router.Group("asnInfo")
	asnInfoRouterWithoutAuth := PublicRouter.Group("asnInfo")

	var asnInfoApi = v1.ApiGroupApp.CfscanApiGroup.AsnInfoApi
	{
		asnInfoRouter.POST("createAsnInfo", asnInfoApi.CreateAsnInfo)             // 新建asnInfo表
		asnInfoRouter.DELETE("deleteAsnInfo", asnInfoApi.DeleteAsnInfo)           // 删除asnInfo表
		asnInfoRouter.DELETE("deleteAsnInfoByIds", asnInfoApi.DeleteAsnInfoByIds) // 批量删除asnInfo表
		asnInfoRouter.PUT("updateAsnInfo", asnInfoApi.UpdateAsnInfo)              // 更新asnInfo表
	}
	{
		asnInfoRouterWithoutRecord.GET("findAsnInfo", asnInfoApi.FindAsnInfo)       // 根据ID获取asnInfo表
		asnInfoRouterWithoutRecord.GET("getAsnInfoList", asnInfoApi.GetAsnInfoList) // 获取asnInfo表列表
	}
	{
		asnInfoRouterWithoutAuth.GET("getAsnInfoPublic", asnInfoApi.GetAsnInfoPublic) // 获取asnInfo表列表
	}
}
