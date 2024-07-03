// 自动生成模板AsnInfo
package cfscan

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// asnInfo表 结构体  AsnInfo
type AsnInfo struct {
	global.GVA_MODEL
	AsnName           string     `json:"asnName" form:"asnName" gorm:"uniqueIndex;column:asn_name;comment:ASN名称;" binding:"required"` //ASN名称
	FullName          string     `json:"fullName" form:"fullName" gorm:"column:full_name;comment:ASN全名;" binding:"required"`          //ASN全名
	Ipv4Counts        *int       `json:"ipv4Counts" form:"ipv4Counts" gorm:"column:ipv4_counts;comment:IPV4数量;"`                      //IPV4数量
	Ipv6Counts        *int       `json:"ipv6Counts" form:"ipv6Counts" gorm:"column:ipv6_counts;comment:IPV6数量;"`                      //IPV6数量
	PeersCounts       *int       `json:"peersCounts" form:"peersCounts" gorm:"column:peers_counts;comment:节点数量;"`                     //节点数量
	Ipv4Peers         *int       `json:"ipv4Peers" form:"ipv4Peers" gorm:"column:ipv4_peers;comment:IPV4节点数量;"`                       //IPV4节点数量
	Ipv6Peers         *int       `json:"ipv6Peers" form:"ipv6Peers" gorm:"column:ipv6_peers;comment:IPV6节点数量;"`                       //IPV6节点数量
	PrefixesCounts    *int       `json:"prefixesCounts" form:"prefixesCounts" gorm:"column:prefixes_counts;comment:IP前缀数量;"`          //IP前缀数量
	Ipv4Prefixies     *int       `json:"ipv4Prefixies" form:"ipv4Prefixies" gorm:"column:ipv4_prefixies;comment:IPV4前缀数量;"`           //IPV4前缀数量
	Ipv6Prefixies     *int       `json:"ipv6Prefixies" form:"ipv6Prefixies" gorm:"column:ipv6_prefixies;comment:IPV6前缀数量;"`           //IPV6前缀数量
	RegionalRegistry  string     `json:"regionalRegistry" form:"regionalRegistry" gorm:"column:regional_registry;comment:地区登记;"`      //地区登记
	TrafficBandwidth  string     `json:"trafficBandwidth" form:"trafficBandwidth" gorm:"column:traffic_bandwidth;comment:带宽估算;"`      //带宽估算
	AllocationStatus  string     `json:"allocationStatus" form:"allocationStatus" gorm:"column:allocation_status;comment:分配状态;"`      //分配状态
	TrafficRatio      string     `json:"trafficRatio" form:"trafficRatio" gorm:"column:traffic_ratio;comment:流量比率;"`                  //流量比率
	AllocationDate    string     `json:"allocationDate" form:"allocationDate" gorm:"column:allocation_date;comment:分配日期;"`            //分配日期
	Website           string     `json:"website" form:"website" gorm:"column:website;comment:官方网址;"`                                  //官方网址
	AllocationCountry string     `json:"allocationCountry" form:"allocationCountry" gorm:"column:allocation_country;comment:分配国家;"`   //分配国家
	Ipv4CIDR          string     `json:"ipv4CIDR" form:"ipv4CIDR" gorm:"column:ipv4CIDR;comment:ipv4 CIDR数据;"`                        //IPV4 CIDR
	Enable            *int       `json:"enable" form:"enable" gorm:"default:1;column:enable;comment:是否开启;" binding:"required"`        //是否开启
	LastCIDRUpdate    *time.Time `json:"lastCIDRUpdate" form:"lastCIDRUpdate" gorm:"column:last_cidr_update;comment:CIDR最后更新时间;"`     //CIDR最后更新时间
}

// TableName asnInfo表 AsnInfo自定义表名 asn_info
func (AsnInfo) TableName() string {
	return "asn_info"
}
