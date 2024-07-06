package cfscan

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ScheduleTaskRouter struct{}

// InitScheduleTaskRouter 初始化 scheduleTask表 路由信息
func (s *ScheduleTaskRouter) InitScheduleTaskRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	scheduleTaskRouter := Router.Group("scheduleTask").Use(middleware.OperationRecord())
	scheduleTaskRouterWithoutRecord := Router.Group("scheduleTask")
	scheduleTaskRouterWithoutAuth := PublicRouter.Group("scheduleTask")

	var scheduleTaskApi = v1.ApiGroupApp.CfscanApiGroup.ScheduleTaskApi
	{
		scheduleTaskRouter.POST("createScheduleTask", scheduleTaskApi.CreateScheduleTask)             // 新建scheduleTask表
		scheduleTaskRouter.DELETE("deleteScheduleTask", scheduleTaskApi.DeleteScheduleTask)           // 删除scheduleTask表
		scheduleTaskRouter.DELETE("deleteScheduleTaskByIds", scheduleTaskApi.DeleteScheduleTaskByIds) // 批量删除scheduleTask表
		scheduleTaskRouter.PUT("updateScheduleTask", scheduleTaskApi.UpdateScheduleTask)              // 更新scheduleTask表
	}
	{
		scheduleTaskRouterWithoutRecord.GET("findScheduleTask", scheduleTaskApi.FindScheduleTask)       // 根据ID获取scheduleTask表
		scheduleTaskRouterWithoutRecord.GET("getScheduleTaskList", scheduleTaskApi.GetScheduleTaskList) // 获取scheduleTask表列表
	}
	{
		scheduleTaskRouterWithoutAuth.GET("getScheduleTaskPublic", scheduleTaskApi.GetScheduleTaskPublic) // 获取scheduleTask表列表
	}
}
