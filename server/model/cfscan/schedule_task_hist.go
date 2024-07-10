// 自动生成模板ScheduleTaskHist
package cfscan

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// scheduleTaskHist表 结构体  ScheduleTaskHist
type ScheduleTaskHist struct {
	global.GVA_MODEL
	ASNName        string    `json:"asnName" form:"asnName" gorm:"column:asn_name;comment:;"`                             //ASN名称
	ScheduleTaskId int       `json:"scheduleTaskId" form:"scheduleTaskId" gorm:"index;column:schedule_task_id;comment:;"` //定时任务ID
	StartTime      time.Time `json:"startTime" form:"startTime" gorm:"column:start_time;comment:;"`                       //起始时间
	EndTime        time.Time `json:"endTime" form:"endTime" gorm:"column:end_time;comment:;"`                             //结束时间
	CostTime       int       `json:"costTime" form:"costTime" gorm:"column:cost_time;comment:;"`                          //耗时
	HistStatus     string    `json:"histStatus" form:"histStatus" gorm:"column:hist_status;comment:;"`                    //任务状态
	TaskResult     string    `json:"taskResult" form:"taskResult" gorm:"column:task_result;comment:;"`                    //任务结果
}

// TableName scheduleTaskHist表 ScheduleTaskHist自定义表名 schedule_task_hist
func (ScheduleTaskHist) TableName() string {
	return "schedule_task_hist"
}
