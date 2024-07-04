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

type SubmitScanApi struct{}

var submitScanService = service.ServiceGroupApp.CfscanServiceGroup.SubmitScanService

// CreateSubmitScan 创建submitScan表
// @Tags SubmitScan
// @Summary 创建submitScan表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cfscan.SubmitScan true "创建submitScan表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /submitScan/createSubmitScan [post]
func (submitScanApi *SubmitScanApi) CreateSubmitScan(c *gin.Context) {
	var submitScan cfscan.SubmitScan
	err := c.ShouldBindJSON(&submitScan)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := submitScanService.CreateSubmitScan(&submitScan); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteSubmitScan 删除submitScan表
// @Tags SubmitScan
// @Summary 删除submitScan表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cfscan.SubmitScan true "删除submitScan表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /submitScan/deleteSubmitScan [delete]
func (submitScanApi *SubmitScanApi) DeleteSubmitScan(c *gin.Context) {
	ID := c.Query("ID")
	if err := submitScanService.DeleteSubmitScan(ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteSubmitScanByIds 批量删除submitScan表
// @Tags SubmitScan
// @Summary 批量删除submitScan表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /submitScan/deleteSubmitScanByIds [delete]
func (submitScanApi *SubmitScanApi) DeleteSubmitScanByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	if err := submitScanService.DeleteSubmitScanByIds(IDs); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateSubmitScan 更新submitScan表
// @Tags SubmitScan
// @Summary 更新submitScan表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cfscan.SubmitScan true "更新submitScan表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /submitScan/updateSubmitScan [put]
func (submitScanApi *SubmitScanApi) UpdateSubmitScan(c *gin.Context) {
	var submitScan cfscan.SubmitScan
	err := c.ShouldBindJSON(&submitScan)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := submitScanService.UpdateSubmitScan(submitScan); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindSubmitScan 用id查询submitScan表
// @Tags SubmitScan
// @Summary 用id查询submitScan表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query cfscan.SubmitScan true "用id查询submitScan表"
// @Success 200 {object} response.Response{data=object{resubmitScan=cfscan.SubmitScan},msg=string} "查询成功"
// @Router /submitScan/findSubmitScan [get]
func (submitScanApi *SubmitScanApi) FindSubmitScan(c *gin.Context) {
	ID := c.Query("ID")
	if resubmitScan, err := submitScanService.GetSubmitScan(ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(resubmitScan, c)
	}
}

// GetSubmitScanList 分页获取submitScan表列表
// @Tags SubmitScan
// @Summary 分页获取submitScan表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query cfscanReq.SubmitScanSearch true "分页获取submitScan表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /submitScan/getSubmitScanList [get]
func (submitScanApi *SubmitScanApi) GetSubmitScanList(c *gin.Context) {
	var pageInfo cfscanReq.SubmitScanSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := submitScanService.GetSubmitScanInfoList(pageInfo); err != nil {
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

// GetSubmitScanPublic 不需要鉴权的submitScan表接口
// @Tags SubmitScan
// @Summary 不需要鉴权的submitScan表接口
// @accept application/json
// @Produce application/json
// @Param data query cfscanReq.SubmitScanSearch true "分页获取submitScan表列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /submitScan/getSubmitScanPublic [get]
func (submitScanApi *SubmitScanApi) GetSubmitScanPublic(c *gin.Context) {
	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的submitScan表接口信息",
	}, "获取成功", c)
}
