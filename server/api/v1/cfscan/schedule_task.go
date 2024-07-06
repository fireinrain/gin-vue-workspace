package cfscan

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cfscan"
	cfscanReq "github.com/flipped-aurora/gin-vue-admin/server/model/cfscan/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ScheduleTaskApi struct{}

var scheduleTaskService = service.ServiceGroupApp.CfscanServiceGroup.ScheduleTaskService

// CreateScheduleTask 创建scheduleTask表
// @Tags ScheduleTask
// @Summary 创建scheduleTask表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cfscan.ScheduleTask true "创建scheduleTask表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /scheduleTask/createScheduleTask [post]
func (scheduleTaskApi *ScheduleTaskApi) CreateScheduleTask(c *gin.Context) {
	var scheduleTask cfscan.ScheduleTask
	err := c.ShouldBindJSON(&scheduleTask)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := scheduleTaskService.CreateScheduleTask(&scheduleTask); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteScheduleTask 删除scheduleTask表
// @Tags ScheduleTask
// @Summary 删除scheduleTask表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cfscan.ScheduleTask true "删除scheduleTask表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /scheduleTask/deleteScheduleTask [delete]
func (scheduleTaskApi *ScheduleTaskApi) DeleteScheduleTask(c *gin.Context) {
	ID := c.Query("ID")
	if err := scheduleTaskService.DeleteScheduleTask(ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteScheduleTaskByIds 批量删除scheduleTask表
// @Tags ScheduleTask
// @Summary 批量删除scheduleTask表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /scheduleTask/deleteScheduleTaskByIds [delete]
func (scheduleTaskApi *ScheduleTaskApi) DeleteScheduleTaskByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	if err := scheduleTaskService.DeleteScheduleTaskByIds(IDs); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateScheduleTask 更新scheduleTask表
// @Tags ScheduleTask
// @Summary 更新scheduleTask表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cfscan.ScheduleTask true "更新scheduleTask表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /scheduleTask/updateScheduleTask [put]
func (scheduleTaskApi *ScheduleTaskApi) UpdateScheduleTask(c *gin.Context) {
	var scheduleTask cfscan.ScheduleTask
	err := c.ShouldBindJSON(&scheduleTask)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := scheduleTaskService.UpdateScheduleTask(scheduleTask); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindScheduleTask 用id查询scheduleTask表
// @Tags ScheduleTask
// @Summary 用id查询scheduleTask表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query cfscan.ScheduleTask true "用id查询scheduleTask表"
// @Success 200 {object} response.Response{data=object{rescheduleTask=cfscan.ScheduleTask},msg=string} "查询成功"
// @Router /scheduleTask/findScheduleTask [get]
func (scheduleTaskApi *ScheduleTaskApi) FindScheduleTask(c *gin.Context) {
	ID := c.Query("ID")
	if rescheduleTask, err := scheduleTaskService.GetScheduleTask(ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(rescheduleTask, c)
	}
}

// GetScheduleTaskList 分页获取scheduleTask表列表
// @Tags ScheduleTask
// @Summary 分页获取scheduleTask表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query cfscanReq.ScheduleTaskSearch true "分页获取scheduleTask表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /scheduleTask/getScheduleTaskList [get]
func (scheduleTaskApi *ScheduleTaskApi) GetScheduleTaskList(c *gin.Context) {
	var pageInfo cfscanReq.ScheduleTaskSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := scheduleTaskService.GetScheduleTaskInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// GetScheduleTaskPublic 不需要鉴权的scheduleTask表接口
// @Tags ScheduleTask
// @Summary 不需要鉴权的scheduleTask表接口
// @accept application/json
// @Produce application/json
// @Param data query cfscanReq.ScheduleTaskSearch true "分页获取scheduleTask表列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /scheduleTask/getScheduleTaskPublic [get]
func (scheduleTaskApi *ScheduleTaskApi) GetScheduleTaskPublic(c *gin.Context) {
	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的scheduleTask表接口信息",
	}, "获取成功", c)
}
