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

type AliveProxyIpsApi struct{}

var aliveProxyIpsService = service.ServiceGroupApp.CfscanServiceGroup.AliveProxyIpsService

// CreateAliveProxyIps 创建aliveProxyIps表
// @Tags AliveProxyIps
// @Summary 创建aliveProxyIps表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cfscan.AliveProxyIps true "创建aliveProxyIps表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /aliveProxyIps/createAliveProxyIps [post]
func (aliveProxyIpsApi *AliveProxyIpsApi) CreateAliveProxyIps(c *gin.Context) {
	var aliveProxyIps cfscan.AliveProxyIps
	err := c.ShouldBindJSON(&aliveProxyIps)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := aliveProxyIpsService.CreateAliveProxyIps(&aliveProxyIps); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteAliveProxyIps 删除aliveProxyIps表
// @Tags AliveProxyIps
// @Summary 删除aliveProxyIps表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cfscan.AliveProxyIps true "删除aliveProxyIps表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /aliveProxyIps/deleteAliveProxyIps [delete]
func (aliveProxyIpsApi *AliveProxyIpsApi) DeleteAliveProxyIps(c *gin.Context) {
	ID := c.Query("ID")
	if err := aliveProxyIpsService.DeleteAliveProxyIps(ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteAliveProxyIpsByIds 批量删除aliveProxyIps表
// @Tags AliveProxyIps
// @Summary 批量删除aliveProxyIps表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /aliveProxyIps/deleteAliveProxyIpsByIds [delete]
func (aliveProxyIpsApi *AliveProxyIpsApi) DeleteAliveProxyIpsByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	if err := aliveProxyIpsService.DeleteAliveProxyIpsByIds(IDs); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateAliveProxyIps 更新aliveProxyIps表
// @Tags AliveProxyIps
// @Summary 更新aliveProxyIps表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cfscan.AliveProxyIps true "更新aliveProxyIps表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /aliveProxyIps/updateAliveProxyIps [put]
func (aliveProxyIpsApi *AliveProxyIpsApi) UpdateAliveProxyIps(c *gin.Context) {
	var aliveProxyIps cfscan.AliveProxyIps
	err := c.ShouldBindJSON(&aliveProxyIps)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := aliveProxyIpsService.UpdateAliveProxyIps(aliveProxyIps); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindAliveProxyIps 用id查询aliveProxyIps表
// @Tags AliveProxyIps
// @Summary 用id查询aliveProxyIps表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query cfscan.AliveProxyIps true "用id查询aliveProxyIps表"
// @Success 200 {object} response.Response{data=object{realiveProxyIps=cfscan.AliveProxyIps},msg=string} "查询成功"
// @Router /aliveProxyIps/findAliveProxyIps [get]
func (aliveProxyIpsApi *AliveProxyIpsApi) FindAliveProxyIps(c *gin.Context) {
	ID := c.Query("ID")
	if realiveProxyIps, err := aliveProxyIpsService.GetAliveProxyIps(ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(realiveProxyIps, c)
	}
}

// GetAliveProxyIpsList 分页获取aliveProxyIps表列表
// @Tags AliveProxyIps
// @Summary 分页获取aliveProxyIps表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query cfscanReq.AliveProxyIpsSearch true "分页获取aliveProxyIps表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /aliveProxyIps/getAliveProxyIpsList [get]
func (aliveProxyIpsApi *AliveProxyIpsApi) GetAliveProxyIpsList(c *gin.Context) {
	var pageInfo cfscanReq.AliveProxyIpsSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := aliveProxyIpsService.GetAliveProxyIpsInfoList(pageInfo); err != nil {
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

// GetAliveProxyIpsPublic 不需要鉴权的aliveProxyIps表接口
// @Tags AliveProxyIps
// @Summary 不需要鉴权的aliveProxyIps表接口
// @accept application/json
// @Produce application/json
// @Param data query cfscanReq.AliveProxyIpsSearch true "分页获取aliveProxyIps表列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /aliveProxyIps/getAliveProxyIpsPublic [get]
func (aliveProxyIpsApi *AliveProxyIpsApi) GetAliveProxyIpsPublic(c *gin.Context) {
	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的aliveProxyIps表接口信息",
	}, "获取成功", c)
}
