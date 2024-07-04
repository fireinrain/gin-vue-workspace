import service from '@/utils/request'

// @Tags SubmitScan
// @Summary 创建submitScan表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SubmitScan true "创建submitScan表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /submitScan/createSubmitScan [post]
export const createSubmitScan = (data) => {
  return service({
    url: '/submitScan/createSubmitScan',
    method: 'post',
    data
  })
}

// @Tags SubmitScan
// @Summary 删除submitScan表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SubmitScan true "删除submitScan表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /submitScan/deleteSubmitScan [delete]
export const deleteSubmitScan = (params) => {
  return service({
    url: '/submitScan/deleteSubmitScan',
    method: 'delete',
    params
  })
}

// @Tags SubmitScan
// @Summary 批量删除submitScan表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除submitScan表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /submitScan/deleteSubmitScan [delete]
export const deleteSubmitScanByIds = (params) => {
  return service({
    url: '/submitScan/deleteSubmitScanByIds',
    method: 'delete',
    params
  })
}

// @Tags SubmitScan
// @Summary 更新submitScan表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SubmitScan true "更新submitScan表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /submitScan/updateSubmitScan [put]
export const updateSubmitScan = (data) => {
  return service({
    url: '/submitScan/updateSubmitScan',
    method: 'put',
    data
  })
}

// @Tags SubmitScan
// @Summary 用id查询submitScan表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.SubmitScan true "用id查询submitScan表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /submitScan/findSubmitScan [get]
export const findSubmitScan = (params) => {
  return service({
    url: '/submitScan/findSubmitScan',
    method: 'get',
    params
  })
}

// @Tags SubmitScan
// @Summary 分页获取submitScan表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取submitScan表列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /submitScan/getSubmitScanList [get]
export const getSubmitScanList = (params) => {
  return service({
    url: '/submitScan/getSubmitScanList',
    method: 'get',
    params
  })
}
