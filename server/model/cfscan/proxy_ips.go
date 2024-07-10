// 自动生成模板ProxyIps
package cfscan

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// proxyIps表 结构体  ProxyIps
type ProxyIps struct {
	global.GVA_MODEL
	AsnNumber     string `json:"asnNumber" form:"asnNumber" gorm:"column:asn_number;comment:;"`             //ASN编号
	Ip            string `json:"ip" form:"ip" gorm:"index;column:ip;comment:;"`                             //IP地址
	Port          int    `json:"port" form:"port" gorm:"index;column:port;comment:;"`                       //端口号
	EnableTls     string `json:"enableTls" form:"enableTls" gorm:"column:enable_tls;comment:;"`             //开启TLS
	DataCenter    string `json:"dataCenter" form:"dataCenter" gorm:"column:data_center;comment:;"`          //数据中心
	Region        string `json:"region" form:"region" gorm:"column:region;comment:;"`                       //地区
	City          string `json:"city" form:"city" gorm:"column:city;comment:;"`                             //城市
	Latency       string `json:"latency" form:"latency" gorm:"column:latency;comment:;"`                    //延迟
	TcpDuration   string `json:"tcpDuration" form:"tcpDuration" gorm:"column:tcp_duration;comment:;"`       //TCP延迟
	DownloadSpeed string `json:"downloadSpeed" form:"downloadSpeed" gorm:"column:download_speed;comment:;"` //下载速度
}

// TableName proxyIps表 ProxyIps自定义表名 proxy_ips
func (ProxyIps) TableName() string {
	return "proxy_ips"
}
