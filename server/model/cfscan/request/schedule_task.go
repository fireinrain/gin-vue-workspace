package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type ScheduleTaskSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	TaskDesc       string     `json:"taskDesc" form:"taskDesc" `
	AsnNumber      string     `json:"asnNumber" form:"asnNumber" `
	AsnDesc        string     `json:"asnDesc" form:"asnDesc" `
	request.PageInfo
	Sort  string `json:"sort" form:"sort"`
	Order string `json:"order" form:"order"`
}
