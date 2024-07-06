package cfscan

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cfscan"
	cfscanReq "github.com/flipped-aurora/gin-vue-admin/server/model/cfscan/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AsnInfoApi struct{}

var asnInfoService = service.ServiceGroupApp.CfscanServiceGroup.AsnInfoService

// CreateAsnInfo 创建asnInfo表
// @Tags AsnInfo
// @Summary 创建asnInfo表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cfscan.AsnInfo true "创建asnInfo表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /asnInfo/createAsnInfo [post]
func (asnInfoApi *AsnInfoApi) CreateAsnInfo(c *gin.Context) {
	var asnInfo cfscan.AsnInfo
	err := c.ShouldBindJSON(&asnInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	asn := asnInfoService.GetASNDetailByASN(&asnInfo)
	if asn == nil {
		err := errors.New("ASN编号尚未被注册")
		response.FailWithMessage(err.Error(), c)
		return
	}
	asn.AsnName = asnInfo.AsnName

	if err := asnInfoService.CreateAsnInfo(asn); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteAsnInfo 删除asnInfo表
// @Tags AsnInfo
// @Summary 删除asnInfo表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cfscan.AsnInfo true "删除asnInfo表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /asnInfo/deleteAsnInfo [delete]
func (asnInfoApi *AsnInfoApi) DeleteAsnInfo(c *gin.Context) {
	ID := c.Query("ID")
	if err := asnInfoService.DeleteAsnInfo(ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteAsnInfoByIds 批量删除asnInfo表
// @Tags AsnInfo
// @Summary 批量删除asnInfo表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /asnInfo/deleteAsnInfoByIds [delete]
func (asnInfoApi *AsnInfoApi) DeleteAsnInfoByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	if err := asnInfoService.DeleteAsnInfoByIds(IDs); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateAsnInfo 更新asnInfo表
// @Tags AsnInfo
// @Summary 更新asnInfo表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cfscan.AsnInfo true "更新asnInfo表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /asnInfo/updateAsnInfo [put]
func (asnInfoApi *AsnInfoApi) UpdateAsnInfo(c *gin.Context) {
	var asnInfo cfscan.AsnInfo
	err := c.ShouldBindJSON(&asnInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := asnInfoService.UpdateAsnInfo(asnInfo); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindAsnInfo 用id查询asnInfo表
// @Tags AsnInfo
// @Summary 用id查询asnInfo表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query cfscan.AsnInfo true "用id查询asnInfo表"
// @Success 200 {object} response.Response{data=object{reasnInfo=cfscan.AsnInfo},msg=string} "查询成功"
// @Router /asnInfo/findAsnInfo [get]
func (asnInfoApi *AsnInfoApi) FindAsnInfo(c *gin.Context) {
	ID := c.Query("ID")
	if reasnInfo, err := asnInfoService.GetAsnInfo(ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(reasnInfo, c)
	}
}

// GetAsnInfoList 分页获取asnInfo表列表
// @Tags AsnInfo
// @Summary 分页获取asnInfo表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query cfscanReq.AsnInfoSearch true "分页获取asnInfo表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /asnInfo/getAsnInfoList [get]
func (asnInfoApi *AsnInfoApi) GetAsnInfoList(c *gin.Context) {
	var pageInfo cfscanReq.AsnInfoSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := asnInfoService.GetAsnInfoInfoList(pageInfo); err != nil {
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

// GetAsnInfoPublic 不需要鉴权的asnInfo表接口
// @Tags AsnInfo
// @Summary 不需要鉴权的asnInfo表接口
// @accept application/json
// @Produce application/json
// @Param data query cfscanReq.AsnInfoSearch true "分页获取asnInfo表列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /asnInfo/getAsnInfoPublic [get]
func (asnInfoApi *AsnInfoApi) GetAsnInfoPublic(c *gin.Context) {
	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的asnInfo表接口信息",
	}, "获取成功", c)
}
