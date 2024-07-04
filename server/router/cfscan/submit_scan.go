package cfscan

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SubmitScanRouter struct{}

// InitSubmitScanRouter 初始化 submitScan表 路由信息
func (s *SubmitScanRouter) InitSubmitScanRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	submitScanRouter := Router.Group("submitScan").Use(middleware.OperationRecord())
	submitScanRouterWithoutRecord := Router.Group("submitScan")
	submitScanRouterWithoutAuth := PublicRouter.Group("submitScan")

	var submitScanApi = v1.ApiGroupApp.CfscanApiGroup.SubmitScanApi
	{
		submitScanRouter.POST("createSubmitScan", submitScanApi.CreateSubmitScan)             // 新建submitScan表
		submitScanRouter.DELETE("deleteSubmitScan", submitScanApi.DeleteSubmitScan)           // 删除submitScan表
		submitScanRouter.DELETE("deleteSubmitScanByIds", submitScanApi.DeleteSubmitScanByIds) // 批量删除submitScan表
		submitScanRouter.PUT("updateSubmitScan", submitScanApi.UpdateSubmitScan)              // 更新submitScan表
	}
	{
		submitScanRouterWithoutRecord.GET("findSubmitScan", submitScanApi.FindSubmitScan)       // 根据ID获取submitScan表
		submitScanRouterWithoutRecord.GET("getSubmitScanList", submitScanApi.GetSubmitScanList) // 获取submitScan表列表
	}
	{
		submitScanRouterWithoutAuth.GET("getSubmitScanPublic", submitScanApi.GetSubmitScanPublic) // 获取submitScan表列表
	}
}
