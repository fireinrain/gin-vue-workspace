package cfscan

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AliveProxyIpsRouter struct{}

// InitAliveProxyIpsRouter 初始化 aliveProxyIps表 路由信息
func (s *AliveProxyIpsRouter) InitAliveProxyIpsRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	aliveProxyIpsRouter := Router.Group("aliveProxyIps").Use(middleware.OperationRecord())
	aliveProxyIpsRouterWithoutRecord := Router.Group("aliveProxyIps")
	aliveProxyIpsRouterWithoutAuth := PublicRouter.Group("aliveProxyIps")

	var aliveProxyIpsApi = v1.ApiGroupApp.CfscanApiGroup.AliveProxyIpsApi
	{
		aliveProxyIpsRouter.POST("createAliveProxyIps", aliveProxyIpsApi.CreateAliveProxyIps)             // 新建aliveProxyIps表
		aliveProxyIpsRouter.DELETE("deleteAliveProxyIps", aliveProxyIpsApi.DeleteAliveProxyIps)           // 删除aliveProxyIps表
		aliveProxyIpsRouter.DELETE("deleteAliveProxyIpsByIds", aliveProxyIpsApi.DeleteAliveProxyIpsByIds) // 批量删除aliveProxyIps表
		aliveProxyIpsRouter.PUT("updateAliveProxyIps", aliveProxyIpsApi.UpdateAliveProxyIps)              // 更新aliveProxyIps表
	}
	{
		aliveProxyIpsRouterWithoutRecord.GET("findAliveProxyIps", aliveProxyIpsApi.FindAliveProxyIps)       // 根据ID获取aliveProxyIps表
		aliveProxyIpsRouterWithoutRecord.GET("getAliveProxyIpsList", aliveProxyIpsApi.GetAliveProxyIpsList) // 获取aliveProxyIps表列表
	}
	{
		aliveProxyIpsRouterWithoutAuth.GET("getAliveProxyIpsPublic", aliveProxyIpsApi.GetAliveProxyIpsPublic) // 获取aliveProxyIps表列表
	}
}
