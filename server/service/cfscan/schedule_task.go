package cfscan

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cfscan"
	cfscanReq "github.com/flipped-aurora/gin-vue-admin/server/model/cfscan/request"
)

type ScheduleTaskService struct{}

// CreateScheduleTask 创建scheduleTask表记录
// Author [piexlmax](https://github.com/piexlmax)
func (scheduleTaskService *ScheduleTaskService) CreateScheduleTask(scheduleTask *cfscan.ScheduleTask) (err error) {
	err = global.GVA_DB.Create(scheduleTask).Error
	return err
}

// DeleteScheduleTask 删除scheduleTask表记录
// Author [piexlmax](https://github.com/piexlmax)
func (scheduleTaskService *ScheduleTaskService) DeleteScheduleTask(ID string) (err error) {
	err = global.GVA_DB.Delete(&cfscan.ScheduleTask{}, "id = ?", ID).Error
	return err
}

// DeleteScheduleTaskByIds 批量删除scheduleTask表记录
// Author [piexlmax](https://github.com/piexlmax)
func (scheduleTaskService *ScheduleTaskService) DeleteScheduleTaskByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]cfscan.ScheduleTask{}, "id in ?", IDs).Error
	return err
}

// UpdateScheduleTask 更新scheduleTask表记录
// Author [piexlmax](https://github.com/piexlmax)
func (scheduleTaskService *ScheduleTaskService) UpdateScheduleTask(scheduleTask cfscan.ScheduleTask) (err error) {
	err = global.GVA_DB.Model(&cfscan.ScheduleTask{}).Where("id = ?", scheduleTask.ID).Updates(&scheduleTask).Error
	return err
}

// GetScheduleTask 根据ID获取scheduleTask表记录
// Author [piexlmax](https://github.com/piexlmax)
func (scheduleTaskService *ScheduleTaskService) GetScheduleTask(ID string) (scheduleTask cfscan.ScheduleTask, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&scheduleTask).Error
	return
}

// GetScheduleTaskInfoList 分页获取scheduleTask表记录
// Author [piexlmax](https://github.com/piexlmax)
func (scheduleTaskService *ScheduleTaskService) GetScheduleTaskInfoList(info cfscanReq.ScheduleTaskSearch) (list []cfscan.ScheduleTask, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&cfscan.ScheduleTask{})
	var scheduleTasks []cfscan.ScheduleTask
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.TaskDesc != "" {
		db = db.Where("task_desc LIKE ?", "%"+info.TaskDesc+"%")
	}
	if info.AsnNumber != "" {
		db = db.Where("asn_number LIKE ?", "%"+info.AsnNumber+"%")
	}
	if info.AsnDesc != "" {
		db = db.Where("asn_desc LIKE ?", "%"+info.AsnDesc+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	var OrderStr string
	orderMap := make(map[string]bool)
	orderMap["asn_number"] = true
	if orderMap[info.Sort] {
		OrderStr = info.Sort
		if info.Order == "descending" {
			OrderStr = OrderStr + " desc"
		}
		db = db.Order(OrderStr)
	} else {
		OrderStr = "id desc"
		db = db.Order(OrderStr)
	}

	err = db.Find(&scheduleTasks).Error
	return scheduleTasks, total, err
}
