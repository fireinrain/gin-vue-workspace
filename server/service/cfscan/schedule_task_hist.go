package cfscan

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cfscan"
	cfscanReq "github.com/flipped-aurora/gin-vue-admin/server/model/cfscan/request"
)

type ScheduleTaskHistService struct{}

// CreateScheduleTaskHist 创建scheduleTaskHist表记录
// Author [piexlmax](https://github.com/piexlmax)
func (scheduleTaskHistService *ScheduleTaskHistService) CreateScheduleTaskHist(scheduleTaskHist *cfscan.ScheduleTaskHist) (err error) {
	err = global.GVA_DB.Create(scheduleTaskHist).Error
	return err
}

// DeleteScheduleTaskHist 删除scheduleTaskHist表记录
// Author [piexlmax](https://github.com/piexlmax)
func (scheduleTaskHistService *ScheduleTaskHistService) DeleteScheduleTaskHist(ID string) (err error) {
	err = global.GVA_DB.Delete(&cfscan.ScheduleTaskHist{}, "id = ?", ID).Error
	return err
}

// DeleteScheduleTaskHistByIds 批量删除scheduleTaskHist表记录
// Author [piexlmax](https://github.com/piexlmax)
func (scheduleTaskHistService *ScheduleTaskHistService) DeleteScheduleTaskHistByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]cfscan.ScheduleTaskHist{}, "id in ?", IDs).Error
	return err
}

// UpdateScheduleTaskHist 更新scheduleTaskHist表记录
// Author [piexlmax](https://github.com/piexlmax)
func (scheduleTaskHistService *ScheduleTaskHistService) UpdateScheduleTaskHist(scheduleTaskHist cfscan.ScheduleTaskHist) (err error) {
	err = global.GVA_DB.Model(&cfscan.ScheduleTaskHist{}).Where("id = ?", scheduleTaskHist.ID).Updates(&scheduleTaskHist).Error
	return err
}

// GetScheduleTaskHist 根据ID获取scheduleTaskHist表记录
// Author [piexlmax](https://github.com/piexlmax)
func (scheduleTaskHistService *ScheduleTaskHistService) GetScheduleTaskHist(ID string) (scheduleTaskHist cfscan.ScheduleTaskHist, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&scheduleTaskHist).Error
	return
}

// GetScheduleTaskHistInfoList 分页获取scheduleTaskHist表记录
// Author [piexlmax](https://github.com/piexlmax)
func (scheduleTaskHistService *ScheduleTaskHistService) GetScheduleTaskHistInfoList(info cfscanReq.ScheduleTaskHistSearch) (list []cfscan.ScheduleTaskHist, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&cfscan.ScheduleTaskHist{})
	var scheduleTaskHists []cfscan.ScheduleTaskHist
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.ASNName != "" {
		db = db.Where("asn_name LIKE ?", "%"+info.ASNName+"%")
	}
	if info.ScheduleTaskId != nil {
		db = db.Where("schedule_task_id = ?", info.ScheduleTaskId)
	}
	if info.HistStatus != "" {
		db = db.Where("hist_status = ?", info.HistStatus)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var OrderStr string
	orderMap := make(map[string]bool)
	orderMap["start_time"] = true
	orderMap["cost_time"] = true
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

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&scheduleTaskHists).Error
	return scheduleTaskHists, total, err
}
