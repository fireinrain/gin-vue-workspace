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

type ScheduleTaskHistApi struct{}

var scheduleTaskHistService = service.ServiceGroupApp.CfscanServiceGroup.ScheduleTaskHistService

// CreateScheduleTaskHist 创建scheduleTaskHist表
// @Tags ScheduleTaskHist
// @Summary 创建scheduleTaskHist表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cfscan.ScheduleTaskHist true "创建scheduleTaskHist表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /scheduleTaskHist/createScheduleTaskHist [post]
func (scheduleTaskHistApi *ScheduleTaskHistApi) CreateScheduleTaskHist(c *gin.Context) {
	var scheduleTaskHist cfscan.ScheduleTaskHist
	err := c.ShouldBindJSON(&scheduleTaskHist)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := scheduleTaskHistService.CreateScheduleTaskHist(&scheduleTaskHist); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteScheduleTaskHist 删除scheduleTaskHist表
// @Tags ScheduleTaskHist
// @Summary 删除scheduleTaskHist表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cfscan.ScheduleTaskHist true "删除scheduleTaskHist表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /scheduleTaskHist/deleteScheduleTaskHist [delete]
func (scheduleTaskHistApi *ScheduleTaskHistApi) DeleteScheduleTaskHist(c *gin.Context) {
	ID := c.Query("ID")
	if err := scheduleTaskHistService.DeleteScheduleTaskHist(ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteScheduleTaskHistByIds 批量删除scheduleTaskHist表
// @Tags ScheduleTaskHist
// @Summary 批量删除scheduleTaskHist表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /scheduleTaskHist/deleteScheduleTaskHistByIds [delete]
func (scheduleTaskHistApi *ScheduleTaskHistApi) DeleteScheduleTaskHistByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	if err := scheduleTaskHistService.DeleteScheduleTaskHistByIds(IDs); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateScheduleTaskHist 更新scheduleTaskHist表
// @Tags ScheduleTaskHist
// @Summary 更新scheduleTaskHist表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cfscan.ScheduleTaskHist true "更新scheduleTaskHist表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /scheduleTaskHist/updateScheduleTaskHist [put]
func (scheduleTaskHistApi *ScheduleTaskHistApi) UpdateScheduleTaskHist(c *gin.Context) {
	var scheduleTaskHist cfscan.ScheduleTaskHist
	err := c.ShouldBindJSON(&scheduleTaskHist)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := scheduleTaskHistService.UpdateScheduleTaskHist(scheduleTaskHist); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindScheduleTaskHist 用id查询scheduleTaskHist表
// @Tags ScheduleTaskHist
// @Summary 用id查询scheduleTaskHist表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query cfscan.ScheduleTaskHist true "用id查询scheduleTaskHist表"
// @Success 200 {object} response.Response{data=object{rescheduleTaskHist=cfscan.ScheduleTaskHist},msg=string} "查询成功"
// @Router /scheduleTaskHist/findScheduleTaskHist [get]
func (scheduleTaskHistApi *ScheduleTaskHistApi) FindScheduleTaskHist(c *gin.Context) {
	ID := c.Query("ID")
	if rescheduleTaskHist, err := scheduleTaskHistService.GetScheduleTaskHist(ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(rescheduleTaskHist, c)
	}
}

// GetScheduleTaskHistList 分页获取scheduleTaskHist表列表
// @Tags ScheduleTaskHist
// @Summary 分页获取scheduleTaskHist表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query cfscanReq.ScheduleTaskHistSearch true "分页获取scheduleTaskHist表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /scheduleTaskHist/getScheduleTaskHistList [get]
func (scheduleTaskHistApi *ScheduleTaskHistApi) GetScheduleTaskHistList(c *gin.Context) {
	var pageInfo cfscanReq.ScheduleTaskHistSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := scheduleTaskHistService.GetScheduleTaskHistInfoList(pageInfo); err != nil {
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

// GetScheduleTaskHistPublic 不需要鉴权的scheduleTaskHist表接口
// @Tags ScheduleTaskHist
// @Summary 不需要鉴权的scheduleTaskHist表接口
// @accept application/json
// @Produce application/json
// @Param data query cfscanReq.ScheduleTaskHistSearch true "分页获取scheduleTaskHist表列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /scheduleTaskHist/getScheduleTaskHistPublic [get]
func (scheduleTaskHistApi *ScheduleTaskHistApi) GetScheduleTaskHistPublic(c *gin.Context) {
	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的scheduleTaskHist表接口信息",
	}, "获取成功", c)
}
