package cfscan

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cfscan"
	cfscanReq "github.com/flipped-aurora/gin-vue-admin/server/model/cfscan/request"
)

type AliveProxyIpsService struct{}

// CreateAliveProxyIps 创建aliveProxyIps表记录
// Author [piexlmax](https://github.com/piexlmax)
func (aliveProxyIpsService *AliveProxyIpsService) CreateAliveProxyIps(aliveProxyIps *cfscan.AliveProxyIps) (err error) {
	err = global.GVA_DB.Create(aliveProxyIps).Error
	return err
}

// DeleteAliveProxyIps 删除aliveProxyIps表记录
// Author [piexlmax](https://github.com/piexlmax)
func (aliveProxyIpsService *AliveProxyIpsService) DeleteAliveProxyIps(ID string) (err error) {
	err = global.GVA_DB.Delete(&cfscan.AliveProxyIps{}, "id = ?", ID).Error
	return err
}

// DeleteAliveProxyIpsByIds 批量删除aliveProxyIps表记录
// Author [piexlmax](https://github.com/piexlmax)
func (aliveProxyIpsService *AliveProxyIpsService) DeleteAliveProxyIpsByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]cfscan.AliveProxyIps{}, "id in ?", IDs).Error
	return err
}

// UpdateAliveProxyIps 更新aliveProxyIps表记录
// Author [piexlmax](https://github.com/piexlmax)
func (aliveProxyIpsService *AliveProxyIpsService) UpdateAliveProxyIps(aliveProxyIps cfscan.AliveProxyIps) (err error) {
	err = global.GVA_DB.Model(&cfscan.AliveProxyIps{}).Where("id = ?", aliveProxyIps.ID).Updates(&aliveProxyIps).Error
	return err
}

// GetAliveProxyIps 根据ID获取aliveProxyIps表记录
// Author [piexlmax](https://github.com/piexlmax)
func (aliveProxyIpsService *AliveProxyIpsService) GetAliveProxyIps(ID string) (aliveProxyIps cfscan.AliveProxyIps, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&aliveProxyIps).Error
	return
}

// GetAliveProxyIpsInfoList 分页获取aliveProxyIps表记录
// Author [piexlmax](https://github.com/piexlmax)
func (aliveProxyIpsService *AliveProxyIpsService) GetAliveProxyIpsInfoList(info cfscanReq.AliveProxyIpsSearch) (list []cfscan.AliveProxyIps, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&cfscan.AliveProxyIps{})
	var aliveProxyIpss []cfscan.AliveProxyIps
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.AsnNumber != "" {
		db = db.Where("asn_number LIKE ?", "%"+info.AsnNumber+"%")
	}
	if info.Ip != "" {
		db = db.Where("ip LIKE ?", "%"+info.Ip+"%")
	}
	if info.Port != nil {
		db = db.Where("port = ?", info.Port)
	}
	if info.EnableTls != "" {
		db = db.Where("enable_tls = ?", info.EnableTls)
	}
	if info.GeoDistance != nil {
		db = db.Where("geo_distance < ?", info.GeoDistance)
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
	if info.Ttl != nil {
		db = db.Where("ttl < ?", info.Ttl)
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
	orderMap["geo_distance"] = true
	orderMap["data_center"] = true
	orderMap["region"] = true
	orderMap["city"] = true
	orderMap["latency"] = true
	orderMap["download_speed"] = true
	orderMap["ttl"] = true
	if orderMap[info.Sort] {
		OrderStr = info.Sort
		if info.Order == "descending" {
			OrderStr = OrderStr + " desc"
		}
		db = db.Order(OrderStr)
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&aliveProxyIpss).Error
	return aliveProxyIpss, total, err
}
