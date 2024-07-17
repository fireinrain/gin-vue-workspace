import service from '@/utils/request'

// @Tags AliveProxyIps
// @Summary 创建aliveProxyIps表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AliveProxyIps true "创建aliveProxyIps表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /aliveProxyIps/createAliveProxyIps [post]
export const createAliveProxyIps = (data) => {
  return service({
    url: '/aliveProxyIps/createAliveProxyIps',
    method: 'post',
    data
  })
}

// @Tags AliveProxyIps
// @Summary 删除aliveProxyIps表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AliveProxyIps true "删除aliveProxyIps表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /aliveProxyIps/deleteAliveProxyIps [delete]
export const deleteAliveProxyIps = (params) => {
  return service({
    url: '/aliveProxyIps/deleteAliveProxyIps',
    method: 'delete',
    params
  })
}

// @Tags AliveProxyIps
// @Summary 批量删除aliveProxyIps表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除aliveProxyIps表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /aliveProxyIps/deleteAliveProxyIps [delete]
export const deleteAliveProxyIpsByIds = (params) => {
  return service({
    url: '/aliveProxyIps/deleteAliveProxyIpsByIds',
    method: 'delete',
    params
  })
}

// @Tags AliveProxyIps
// @Summary 更新aliveProxyIps表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AliveProxyIps true "更新aliveProxyIps表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /aliveProxyIps/updateAliveProxyIps [put]
export const updateAliveProxyIps = (data) => {
  return service({
    url: '/aliveProxyIps/updateAliveProxyIps',
    method: 'put',
    data
  })
}

// @Tags AliveProxyIps
// @Summary 用id查询aliveProxyIps表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.AliveProxyIps true "用id查询aliveProxyIps表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /aliveProxyIps/findAliveProxyIps [get]
export const findAliveProxyIps = (params) => {
  return service({
    url: '/aliveProxyIps/findAliveProxyIps',
    method: 'get',
    params
  })
}

// @Tags AliveProxyIps
// @Summary 分页获取aliveProxyIps表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取aliveProxyIps表列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /aliveProxyIps/getAliveProxyIpsList [get]
export const getAliveProxyIpsList = (params) => {
  return service({
    url: '/aliveProxyIps/getAliveProxyIpsList',
    method: 'get',
    params
  })
}
