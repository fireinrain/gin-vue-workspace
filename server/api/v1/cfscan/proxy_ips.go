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

type ProxyIpsApi struct{}

var proxyIpsService = service.ServiceGroupApp.CfscanServiceGroup.ProxyIpsService

// CreateProxyIps 创建proxyIps表
// @Tags ProxyIps
// @Summary 创建proxyIps表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cfscan.ProxyIps true "创建proxyIps表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /proxyIps/createProxyIps [post]
func (proxyIpsApi *ProxyIpsApi) CreateProxyIps(c *gin.Context) {
	var proxyIps cfscan.ProxyIps
	err := c.ShouldBindJSON(&proxyIps)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := proxyIpsService.CreateProxyIps(&proxyIps); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteProxyIps 删除proxyIps表
// @Tags ProxyIps
// @Summary 删除proxyIps表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cfscan.ProxyIps true "删除proxyIps表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /proxyIps/deleteProxyIps [delete]
func (proxyIpsApi *ProxyIpsApi) DeleteProxyIps(c *gin.Context) {
	ID := c.Query("ID")
	if err := proxyIpsService.DeleteProxyIps(ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteProxyIpsByIds 批量删除proxyIps表
// @Tags ProxyIps
// @Summary 批量删除proxyIps表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /proxyIps/deleteProxyIpsByIds [delete]
func (proxyIpsApi *ProxyIpsApi) DeleteProxyIpsByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	if err := proxyIpsService.DeleteProxyIpsByIds(IDs); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateProxyIps 更新proxyIps表
// @Tags ProxyIps
// @Summary 更新proxyIps表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cfscan.ProxyIps true "更新proxyIps表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /proxyIps/updateProxyIps [put]
func (proxyIpsApi *ProxyIpsApi) UpdateProxyIps(c *gin.Context) {
	var proxyIps cfscan.ProxyIps
	err := c.ShouldBindJSON(&proxyIps)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := proxyIpsService.UpdateProxyIps(proxyIps); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindProxyIps 用id查询proxyIps表
// @Tags ProxyIps
// @Summary 用id查询proxyIps表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query cfscan.ProxyIps true "用id查询proxyIps表"
// @Success 200 {object} response.Response{data=object{reproxyIps=cfscan.ProxyIps},msg=string} "查询成功"
// @Router /proxyIps/findProxyIps [get]
func (proxyIpsApi *ProxyIpsApi) FindProxyIps(c *gin.Context) {
	ID := c.Query("ID")
	if reproxyIps, err := proxyIpsService.GetProxyIps(ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(reproxyIps, c)
	}
}

// GetProxyIpsList 分页获取proxyIps表列表
// @Tags ProxyIps
// @Summary 分页获取proxyIps表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query cfscanReq.ProxyIpsSearch true "分页获取proxyIps表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /proxyIps/getProxyIpsList [get]
func (proxyIpsApi *ProxyIpsApi) GetProxyIpsList(c *gin.Context) {
	var pageInfo cfscanReq.ProxyIpsSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := proxyIpsService.GetProxyIpsInfoList(pageInfo); err != nil {
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

// GetProxyIpsPublic 不需要鉴权的proxyIps表接口
// @Tags ProxyIps
// @Summary 不需要鉴权的proxyIps表接口
// @accept application/json
// @Produce application/json
// @Param data query cfscanReq.ProxyIpsSearch true "分页获取proxyIps表列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /proxyIps/getProxyIpsPublic [get]
func (proxyIpsApi *ProxyIpsApi) GetProxyIpsPublic(c *gin.Context) {
	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的proxyIps表接口信息",
	}, "获取成功", c)
}
