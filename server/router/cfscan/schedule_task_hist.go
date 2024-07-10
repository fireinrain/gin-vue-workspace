package cfscan

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ScheduleTaskHistRouter struct{}

// InitScheduleTaskHistRouter 初始化 scheduleTaskHist表 路由信息
func (s *ScheduleTaskHistRouter) InitScheduleTaskHistRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	scheduleTaskHistRouter := Router.Group("scheduleTaskHist").Use(middleware.OperationRecord())
	scheduleTaskHistRouterWithoutRecord := Router.Group("scheduleTaskHist")
	scheduleTaskHistRouterWithoutAuth := PublicRouter.Group("scheduleTaskHist")

	var scheduleTaskHistApi = v1.ApiGroupApp.CfscanApiGroup.ScheduleTaskHistApi
	{
		scheduleTaskHistRouter.POST("createScheduleTaskHist", scheduleTaskHistApi.CreateScheduleTaskHist)             // 新建scheduleTaskHist表
		scheduleTaskHistRouter.DELETE("deleteScheduleTaskHist", scheduleTaskHistApi.DeleteScheduleTaskHist)           // 删除scheduleTaskHist表
		scheduleTaskHistRouter.DELETE("deleteScheduleTaskHistByIds", scheduleTaskHistApi.DeleteScheduleTaskHistByIds) // 批量删除scheduleTaskHist表
		scheduleTaskHistRouter.PUT("updateScheduleTaskHist", scheduleTaskHistApi.UpdateScheduleTaskHist)              // 更新scheduleTaskHist表
	}
	{
		scheduleTaskHistRouterWithoutRecord.GET("findScheduleTaskHist", scheduleTaskHistApi.FindScheduleTaskHist)       // 根据ID获取scheduleTaskHist表
		scheduleTaskHistRouterWithoutRecord.GET("getScheduleTaskHistList", scheduleTaskHistApi.GetScheduleTaskHistList) // 获取scheduleTaskHist表列表
	}
	{
		scheduleTaskHistRouterWithoutAuth.GET("getScheduleTaskHistPublic", scheduleTaskHistApi.GetScheduleTaskHistPublic) // 获取scheduleTaskHist表列表
	}
}
