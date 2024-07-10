package cfscan

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ProxyIpsRouter struct{}

// InitProxyIpsRouter 初始化 proxyIps表 路由信息
func (s *ProxyIpsRouter) InitProxyIpsRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	proxyIpsRouter := Router.Group("proxyIps").Use(middleware.OperationRecord())
	proxyIpsRouterWithoutRecord := Router.Group("proxyIps")
	proxyIpsRouterWithoutAuth := PublicRouter.Group("proxyIps")

	var proxyIpsApi = v1.ApiGroupApp.CfscanApiGroup.ProxyIpsApi
	{
		proxyIpsRouter.POST("createProxyIps", proxyIpsApi.CreateProxyIps)             // 新建proxyIps表
		proxyIpsRouter.DELETE("deleteProxyIps", proxyIpsApi.DeleteProxyIps)           // 删除proxyIps表
		proxyIpsRouter.DELETE("deleteProxyIpsByIds", proxyIpsApi.DeleteProxyIpsByIds) // 批量删除proxyIps表
		proxyIpsRouter.PUT("updateProxyIps", proxyIpsApi.UpdateProxyIps)              // 更新proxyIps表
	}
	{
		proxyIpsRouterWithoutRecord.GET("findProxyIps", proxyIpsApi.FindProxyIps)       // 根据ID获取proxyIps表
		proxyIpsRouterWithoutRecord.GET("getProxyIpsList", proxyIpsApi.GetProxyIpsList) // 获取proxyIps表列表
	}
	{
		proxyIpsRouterWithoutAuth.GET("getProxyIpsPublic", proxyIpsApi.GetProxyIpsPublic) // 获取proxyIps表列表
	}
}
