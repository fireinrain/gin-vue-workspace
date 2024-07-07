package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type SubmitScanSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	ScanDesc       string     `json:"scanDesc" form:"scanDesc" `
	ScanType       string     `json:"scanType" form:"scanType" `
	AsnNumber      string     `json:"asnNumber" form:"asnNumber" `
	ScanStatus     string     `json:"scanStatus" form:"scanStatus" `
	request.PageInfo
	Sort  string `json:"sort" form:"sort"`
	Order string `json:"order" form:"order"`
}
