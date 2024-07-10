package cfscan

import "time"

type SpeedTestResultE struct {
	AsnNumber  string `json:"asn_number"`
	IP         string `json:"ip"`          // IP地址
	Port       int    `json:"port"`        // 端口
	EnableTLS  bool   `json:"enable_tls"`  //是否开启TLS
	DataCenter string `json:"data_center"` // 数据中心
	Region     string `json:"region"`      // 地区
	City       string `json:"city"`        // 城市
	//CF下载延迟
	Latency string `json:"latency"` // 延迟
	//CF代理判断延迟
	TcpDuration   time.Duration `json:"tcp_duration"`   // TCP请求延迟
	DownloadSpeed string        `json:"download_speed"` // 下载速度
}
