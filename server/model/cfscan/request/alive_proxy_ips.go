package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type AliveProxyIpsSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	AsnNumber      string     `json:"asnNumber" form:"asnNumber" `
	Ip             string     `json:"ip" form:"ip" `
	Port           *int       `json:"port" form:"port" `
	EnableTls      string     `json:"enableTls" form:"enableTls" `
	GeoDistance    *int       `json:"geoDistance" form:"geoDistance" `
	DataCenter     string     `json:"dataCenter" form:"dataCenter" `
	Region         string     `json:"region" form:"region" `
	City           string     `json:"city" form:"city" `
	Latency        string     `json:"latency" form:"latency" `
	Ttl            *int       `json:"ttl" form:"ttl" `
	request.PageInfo
	Sort  string `json:"sort" form:"sort"`
	Order string `json:"order" form:"order"`
}
