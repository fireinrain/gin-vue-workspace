// 自动生成模板AliveProxyIps
package cfscan

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// aliveProxyIps表 结构体  AliveProxyIps
type AliveProxyIps struct {
	global.GVA_MODEL
	AsnNumber     string `json:"asnNumber" form:"asnNumber" gorm:"column:asn_number;comment:;"`             //ASN编号
	AsnDesc       string `json:"asnDesc" form:"asnDesc" gorm:"column:asn_desc;comment:;"`                   //ASN描述
	Ip            string `json:"ip" form:"ip" gorm:"index;column:ip;comment:;"`                             //IP地址
	Port          *int   `json:"port" form:"port" gorm:"index;column:port;comment:;"`                       //端口
	EnableTls     string `json:"enableTls" form:"enableTls" gorm:"column:enable_tls;comment:;"`             //开启TLS
	GeoDistance   *int   `json:"geoDistance" form:"geoDistance" gorm:"column:geo_distance;comment:;"`       //物理距离
	DataCenter    string `json:"dataCenter" form:"dataCenter" gorm:"column:data_center;comment:;"`          //数据中心
	Region        string `json:"region" form:"region" gorm:"column:region;comment:;"`                       //地区
	City          string `json:"city" form:"city" gorm:"column:city;comment:;"`                             //城市
	Latency       string `json:"latency" form:"latency" gorm:"column:latency;comment:;"`                    //延迟
	TcpDuration   *int   `json:"tcpDuration" form:"tcpDuration" gorm:"column:tcp_duration;comment:;"`       //TCP延迟
	DownloadSpeed string `json:"downloadSpeed" form:"downloadSpeed" gorm:"column:download_speed;comment:;"` //下载速度
	Ttl           *int   `json:"ttl" form:"ttl" gorm:"column:ttl;comment:;"`                                //存活时间
	DescStr       string `json:"descStr" form:"descStr" gorm:"column:desc_str;comment:;"`                   //IP描述
}

// TableName aliveProxyIps表 AliveProxyIps自定义表名 alive_proxy_ips
func (AliveProxyIps) TableName() string {
	return "alive_proxy_ips"
}
