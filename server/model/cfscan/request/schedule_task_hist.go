package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type ScheduleTaskHistSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	ASNName        string     `json:"asnName" form:"asnName" `
	ScheduleTaskId *int       `json:"scheduleTaskId" form:"scheduleTaskId" `
	HistStatus     string     `json:"histStatus" form:"histStatus" `
	request.PageInfo
	Sort  string `json:"sort" form:"sort"`
	Order string `json:"order" form:"order"`
}
