import service from '@/utils/request'

// @Tags ScheduleTask
// @Summary 创建scheduleTask表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ScheduleTask true "创建scheduleTask表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /scheduleTask/createScheduleTask [post]
export const createScheduleTask = (data) => {
  return service({
    url: '/scheduleTask/createScheduleTask',
    method: 'post',
    data
  })
}

// @Tags ScheduleTask
// @Summary 删除scheduleTask表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ScheduleTask true "删除scheduleTask表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /scheduleTask/deleteScheduleTask [delete]
export const deleteScheduleTask = (params) => {
  return service({
    url: '/scheduleTask/deleteScheduleTask',
    method: 'delete',
    params
  })
}

// @Tags ScheduleTask
// @Summary 批量删除scheduleTask表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除scheduleTask表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /scheduleTask/deleteScheduleTask [delete]
export const deleteScheduleTaskByIds = (params) => {
  return service({
    url: '/scheduleTask/deleteScheduleTaskByIds',
    method: 'delete',
    params
  })
}

// @Tags ScheduleTask
// @Summary 更新scheduleTask表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ScheduleTask true "更新scheduleTask表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /scheduleTask/updateScheduleTask [put]
export const updateScheduleTask = (data) => {
  return service({
    url: '/scheduleTask/updateScheduleTask',
    method: 'put',
    data
  })
}

// @Tags ScheduleTask
// @Summary 用id查询scheduleTask表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.ScheduleTask true "用id查询scheduleTask表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /scheduleTask/findScheduleTask [get]
export const findScheduleTask = (params) => {
  return service({
    url: '/scheduleTask/findScheduleTask',
    method: 'get',
    params
  })
}

// @Tags ScheduleTask
// @Summary 分页获取scheduleTask表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取scheduleTask表列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /scheduleTask/getScheduleTaskList [get]
export const getScheduleTaskList = (params) => {
  return service({
    url: '/scheduleTask/getScheduleTaskList',
    method: 'get',
    params
  })
}
