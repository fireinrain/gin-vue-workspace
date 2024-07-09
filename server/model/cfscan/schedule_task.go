// 自动生成模板ScheduleTask
package cfscan

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// scheduleTask表 结构体  ScheduleTask
type ScheduleTask struct {
	global.GVA_MODEL
	TaskDesc   string    `json:"taskDesc" form:"taskDesc" gorm:"column:task_desc;comment:任务描述;" binding:"required"`
	AsnNumber  string    `json:"asnNumber" form:"asnNumber" gorm:"index;column:asn_number;comment:ASN编号;" binding:"required"`
	AsnDesc    string    `json:"asnDesc" form:"asnDesc" gorm:"column:asn_desc;comment:ASN描述;" binding:"required"`
	CrontabStr string    `json:"crontabStr" form:"crontabStr" gorm:"column:crontab_str;comment:定时表达式;" binding:"required"`
	TaskConfig string    `json:"taskConfig" form:"taskConfig" gorm:"column:task_config;comment:任务配置;" binding:"required"`
	Enable     string    `json:"enable" form:"enable" gorm:"default:0;column:enable;comment:是否开启;"`
	TaskStatus string    `json:"taskStatus" form:"taskStatus" gorm:"default:1;column:task_status;comment:任务状态;"`
	LastRunAt  time.Time `json:"lastRunAt" form:"lastRunAt" gorm:"column:last_run_at;comment:上次运行时间;"`
	NextRunAt  time.Time `json:"nextRunAt" form:"nextRunAt" gorm:"column:next_run_at;comment:下次运行时间;"`
}

// TableName scheduleTask表 ScheduleTask自定义表名 schedule_task
func (ScheduleTask) TableName() string {
	return "schedule_task"
}
