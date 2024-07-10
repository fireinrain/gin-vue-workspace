import service from '@/utils/request'

// @Tags ProxyIps
// @Summary 创建proxyIps表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ProxyIps true "创建proxyIps表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /proxyIps/createProxyIps [post]
export const createProxyIps = (data) => {
  return service({
    url: '/proxyIps/createProxyIps',
    method: 'post',
    data
  })
}

// @Tags ProxyIps
// @Summary 删除proxyIps表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ProxyIps true "删除proxyIps表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /proxyIps/deleteProxyIps [delete]
export const deleteProxyIps = (params) => {
  return service({
    url: '/proxyIps/deleteProxyIps',
    method: 'delete',
    params
  })
}

// @Tags ProxyIps
// @Summary 批量删除proxyIps表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除proxyIps表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /proxyIps/deleteProxyIps [delete]
export const deleteProxyIpsByIds = (params) => {
  return service({
    url: '/proxyIps/deleteProxyIpsByIds',
    method: 'delete',
    params
  })
}

// @Tags ProxyIps
// @Summary 更新proxyIps表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ProxyIps true "更新proxyIps表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /proxyIps/updateProxyIps [put]
export const updateProxyIps = (data) => {
  return service({
    url: '/proxyIps/updateProxyIps',
    method: 'put',
    data
  })
}

// @Tags ProxyIps
// @Summary 用id查询proxyIps表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.ProxyIps true "用id查询proxyIps表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /proxyIps/findProxyIps [get]
export const findProxyIps = (params) => {
  return service({
    url: '/proxyIps/findProxyIps',
    method: 'get',
    params
  })
}

// @Tags ProxyIps
// @Summary 分页获取proxyIps表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取proxyIps表列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /proxyIps/getProxyIpsList [get]
export const getProxyIpsList = (params) => {
  return service({
    url: '/proxyIps/getProxyIpsList',
    method: 'get',
    params
  })
}
