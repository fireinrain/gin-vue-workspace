// 自动生成模板SubmitScan
package cfscan

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// submitScan表 结构体  SubmitScan
type SubmitScan struct {
	global.GVA_MODEL
	ScanDesc        string `json:"scanDesc" form:"scanDesc" gorm:"column:scan_desc;comment:;" binding:"required"`                                //scanDesc字段
	ScanType        string `json:"scanType" form:"scanType" gorm:"column:scan_type;comment:;" binding:"required"`                                //scanType字段
	AsnNumber       string `json:"asnNumber" form:"asnNumber" gorm:"column:asn_number;comment:;"`                                                //asnNumber字段
	IpinfoType      string `json:"ipinfoType" form:"ipinfoType" gorm:"column:ipinfo_type;comment:;"`                                             //ipinfoType字段
	IpinfoList      string `json:"ipinfoList" form:"ipinfoList" gorm:"column:ipinfo_list;comment:;"`                                             //ipinfoList字段
	IpbatchSize     *int   `json:"ipbatchSize" form:"ipbatchSize" gorm:"column:ipbatch_size;comment:;"`                                          //ipbatchSize字段
	EnableTls       string `json:"enableTls" form:"enableTls" gorm:"default:1;column:enable_tls;comment:;" binding:"required"`                   //enableTls字段
	ScanPorts       string `json:"scanPorts" form:"scanPorts" gorm:"default:443;column:scan_ports;comment:;" binding:"required"`                 //scanPorts字段
	ScanRate        *int   `json:"scanRate" form:"scanRate" gorm:"default:10000;column:scan_rate;comment:;" binding:"required"`                  //scanRate字段
	IpcheckThread   *int   `json:"ipcheckThread" form:"ipcheckThread" gorm:"default:100;column:ipcheck_thread;comment:;" binding:"required"`     //ipcheckThread字段
	EnableSpeedtest string `json:"enableSpeedtest" form:"enableSpeedtest" gorm:"default:1;column:enable_speedtest;comment:;" binding:"required"` //enableSpeedtest字段
	ScanStatus      string `json:"scanStatus" form:"scanStatus" gorm:"column:scan_status;comment:;"`                                             //scanStatus字段
	ScanResult      string `json:"scanResult" form:"scanResult" gorm:"column:scan_result;comment:;"`                                             //scanResult字段
}

// TableName submitScan表 SubmitScan自定义表名 submit_scan
func (SubmitScan) TableName() string {
	return "submit_scan"
}
