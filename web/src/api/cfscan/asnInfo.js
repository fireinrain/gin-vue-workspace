import service from '@/utils/request'

// @Tags AsnInfo
// @Summary 创建asnInfo表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AsnInfo true "创建asnInfo表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /asnInfo/createAsnInfo [post]
export const createAsnInfo = (data) => {
  return service({
    url: '/asnInfo/createAsnInfo',
    method: 'post',
    data
  })
}

// @Tags AsnInfo
// @Summary 删除asnInfo表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AsnInfo true "删除asnInfo表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /asnInfo/deleteAsnInfo [delete]
export const deleteAsnInfo = (params) => {
  return service({
    url: '/asnInfo/deleteAsnInfo',
    method: 'delete',
    params
  })
}

// @Tags AsnInfo
// @Summary 批量删除asnInfo表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除asnInfo表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /asnInfo/deleteAsnInfo [delete]
export const deleteAsnInfoByIds = (params) => {
  return service({
    url: '/asnInfo/deleteAsnInfoByIds',
    method: 'delete',
    params
  })
}

// @Tags AsnInfo
// @Summary 更新asnInfo表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AsnInfo true "更新asnInfo表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /asnInfo/updateAsnInfo [put]
export const updateAsnInfo = (data) => {
  return service({
    url: '/asnInfo/updateAsnInfo',
    method: 'put',
    data
  })
}

// @Tags AsnInfo
// @Summary 用id查询asnInfo表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.AsnInfo true "用id查询asnInfo表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /asnInfo/findAsnInfo [get]
export const findAsnInfo = (params) => {
  return service({
    url: '/asnInfo/findAsnInfo',
    method: 'get',
    params
  })
}

// @Tags AsnInfo
// @Summary 分页获取asnInfo表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取asnInfo表列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /asnInfo/getAsnInfoList [get]
export const getAsnInfoList = (params) => {
  return service({
    url: '/asnInfo/getAsnInfoList',
    method: 'get',
    params
  })
}
