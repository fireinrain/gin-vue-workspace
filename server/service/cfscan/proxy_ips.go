package cfscan

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cfscan"
	cfscanReq "github.com/flipped-aurora/gin-vue-admin/server/model/cfscan/request"
)

type ProxyIpsService struct{}

// CreateProxyIps 创建proxyIps表记录
// Author [piexlmax](https://github.com/piexlmax)
func (proxyIpsService *ProxyIpsService) CreateProxyIps(proxyIps *cfscan.ProxyIps) (err error) {
	err = global.GVA_DB.Create(proxyIps).Error
	return err
}

// DeleteProxyIps 删除proxyIps表记录
// Author [piexlmax](https://github.com/piexlmax)
func (proxyIpsService *ProxyIpsService) DeleteProxyIps(ID string) (err error) {
	err = global.GVA_DB.Delete(&cfscan.ProxyIps{}, "id = ?", ID).Error
	return err
}

// DeleteProxyIpsByIds 批量删除proxyIps表记录
// Author [piexlmax](https://github.com/piexlmax)
func (proxyIpsService *ProxyIpsService) DeleteProxyIpsByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]cfscan.ProxyIps{}, "id in ?", IDs).Error
	return err
}

// UpdateProxyIps 更新proxyIps表记录
// Author [piexlmax](https://github.com/piexlmax)
func (proxyIpsService *ProxyIpsService) UpdateProxyIps(proxyIps cfscan.ProxyIps) (err error) {
	err = global.GVA_DB.Model(&cfscan.ProxyIps{}).Where("id = ?", proxyIps.ID).Updates(&proxyIps).Error
	return err
}

// GetProxyIps 根据ID获取proxyIps表记录
// Author [piexlmax](https://github.com/piexlmax)
func (proxyIpsService *ProxyIpsService) GetProxyIps(ID string) (proxyIps cfscan.ProxyIps, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&proxyIps).Error
	return
}

// GetProxyIpsInfoList 分页获取proxyIps表记录
// Author [piexlmax](https://github.com/piexlmax)
func (proxyIpsService *ProxyIpsService) GetProxyIpsInfoList(info cfscanReq.ProxyIpsSearch) (list []cfscan.ProxyIps, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&cfscan.ProxyIps{})
	var proxyIpss []cfscan.ProxyIps
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.AsnNumber != "" {
		db = db.Where("asn_number LIKE ?", "%"+info.AsnNumber+"%")
	}
	if info.Ip != "" {
		db = db.Where("ip = ?", info.Ip)
	}
	if info.Port != nil {
		db = db.Where("port = ?", info.Port)
	}
	if info.EnableTls != "" {
		db = db.Where("enable_tls = ?", info.EnableTls)
	}
	if info.DataCenter != "" {
		db = db.Where("data_center LIKE ?", "%"+info.DataCenter+"%")
	}
	if info.Region != "" {
		db = db.Where("region LIKE ?", "%"+info.Region+"%")
	}
	if info.City != "" {
		db = db.Where("city LIKE ?", "%"+info.City+"%")
	}
	if info.Latency != "" {
		db = db.Where("latency < ?", info.Latency)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var OrderStr string
	orderMap := make(map[string]bool)
	orderMap["asn_number"] = true
	orderMap["ip"] = true
	orderMap["port"] = true
	orderMap["enable_tls"] = true
	orderMap["data_center"] = true
	orderMap["region"] = true
	orderMap["city"] = true
	orderMap["latency"] = true
	orderMap["download_speed"] = true
	if orderMap[info.Sort] {
		OrderStr = info.Sort
		if info.Order == "descending" {
			OrderStr = OrderStr + " desc"
		}
		db = db.Order(OrderStr)
	} else {
		OrderStr = "id desc"
		db = db.Order(OrderStr)
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&proxyIpss).Error
	return proxyIpss, total, err
}
