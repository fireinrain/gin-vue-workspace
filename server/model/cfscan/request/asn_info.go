package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/cfscan"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type AsnInfoSearch struct {
	cfscan.AsnInfoParams
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	request.PageInfo
}
