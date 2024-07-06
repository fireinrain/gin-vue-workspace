// 自动生成模板ScheduleTask
package cfscan

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// scheduleTask表 结构体  ScheduleTask
type ScheduleTask struct {
	global.GVA_MODEL
	TaskDesc   string `json:"taskDesc" form:"taskDesc" gorm:"column:task_desc;comment:;" binding:"required"`          //任务描述
	AsnNumber  string `json:"asnNumber" form:"asnNumber" gorm:"index;column:asn_number;comment:;" binding:"required"` //ASN编号
	AsnDesc    string `json:"asnDesc" form:"asnDesc" gorm:"column:asn_desc;comment:;" binding:"required"`             //ASN描述
	CrontabStr string `json:"crontabStr" form:"crontabStr" gorm:"column:crontab_str;comment:;" binding:"required"`    //定时表达式
	TaskConfig string `json:"taskConfig" form:"taskConfig" gorm:"column:task_config;comment:;" binding:"required"`    //任务配置
	Enable     string `json:"enable" form:"enable" gorm:"default:0;column:enable;comment:;"`                          //是否开启
	TaskStatus string `json:"taskStatus" form:"taskStatus" gorm:"default:1;column:task_status;comment:;"`             //任务状态
}

// TableName scheduleTask表 ScheduleTask自定义表名 schedule_task
func (ScheduleTask) TableName() string {
	return "schedule_task"
}
