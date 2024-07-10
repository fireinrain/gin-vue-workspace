import service from '@/utils/request'

// @Tags ScheduleTaskHist
// @Summary 创建scheduleTaskHist表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ScheduleTaskHist true "创建scheduleTaskHist表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /scheduleTaskHist/createScheduleTaskHist [post]
export const createScheduleTaskHist = (data) => {
  return service({
    url: '/scheduleTaskHist/createScheduleTaskHist',
    method: 'post',
    data
  })
}

// @Tags ScheduleTaskHist
// @Summary 删除scheduleTaskHist表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ScheduleTaskHist true "删除scheduleTaskHist表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /scheduleTaskHist/deleteScheduleTaskHist [delete]
export const deleteScheduleTaskHist = (params) => {
  return service({
    url: '/scheduleTaskHist/deleteScheduleTaskHist',
    method: 'delete',
    params
  })
}

// @Tags ScheduleTaskHist
// @Summary 批量删除scheduleTaskHist表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除scheduleTaskHist表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /scheduleTaskHist/deleteScheduleTaskHist [delete]
export const deleteScheduleTaskHistByIds = (params) => {
  return service({
    url: '/scheduleTaskHist/deleteScheduleTaskHistByIds',
    method: 'delete',
    params
  })
}

// @Tags ScheduleTaskHist
// @Summary 更新scheduleTaskHist表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ScheduleTaskHist true "更新scheduleTaskHist表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /scheduleTaskHist/updateScheduleTaskHist [put]
export const updateScheduleTaskHist = (data) => {
  return service({
    url: '/scheduleTaskHist/updateScheduleTaskHist',
    method: 'put',
    data
  })
}

// @Tags ScheduleTaskHist
// @Summary 用id查询scheduleTaskHist表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.ScheduleTaskHist true "用id查询scheduleTaskHist表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /scheduleTaskHist/findScheduleTaskHist [get]
export const findScheduleTaskHist = (params) => {
  return service({
    url: '/scheduleTaskHist/findScheduleTaskHist',
    method: 'get',
    params
  })
}

// @Tags ScheduleTaskHist
// @Summary 分页获取scheduleTaskHist表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取scheduleTaskHist表列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /scheduleTaskHist/getScheduleTaskHistList [get]
export const getScheduleTaskHistList = (params) => {
  return service({
    url: '/scheduleTaskHist/getScheduleTaskHistList',
    method: 'get',
    params
  })
}
