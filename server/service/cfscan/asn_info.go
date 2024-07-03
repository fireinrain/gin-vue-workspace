package cfscan

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cfscan"
	cfscanReq "github.com/flipped-aurora/gin-vue-admin/server/model/cfscan/request"
)

type AsnInfoService struct{}

// CreateAsnInfo 创建asnInfo表记录
// Author [piexlmax](https://github.com/piexlmax)
func (asnInfoService *AsnInfoService) CreateAsnInfo(asnInfo *cfscan.AsnInfo) (err error) {
	err = global.GVA_DB.Create(asnInfo).Error
	return err
}

// DeleteAsnInfo 删除asnInfo表记录
// Author [piexlmax](https://github.com/piexlmax)
func (asnInfoService *AsnInfoService) DeleteAsnInfo(ID string) (err error) {
	err = global.GVA_DB.Delete(&cfscan.AsnInfo{}, "id = ?", ID).Error
	return err
}

// DeleteAsnInfoByIds 批量删除asnInfo表记录
// Author [piexlmax](https://github.com/piexlmax)
func (asnInfoService *AsnInfoService) DeleteAsnInfoByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]cfscan.AsnInfo{}, "id in ?", IDs).Error
	return err
}

// UpdateAsnInfo 更新asnInfo表记录
// Author [piexlmax](https://github.com/piexlmax)
func (asnInfoService *AsnInfoService) UpdateAsnInfo(asnInfo cfscan.AsnInfo) (err error) {
	err = global.GVA_DB.Model(&cfscan.AsnInfo{}).Where("id = ?", asnInfo.ID).Updates(&asnInfo).Error
	return err
}

// GetAsnInfo 根据ID获取asnInfo表记录
// Author [piexlmax](https://github.com/piexlmax)
func (asnInfoService *AsnInfoService) GetAsnInfo(ID string) (asnInfo cfscan.AsnInfo, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&asnInfo).Error
	return
}

// GetAsnInfoInfoList 分页获取asnInfo表记录
// Author [piexlmax](https://github.com/piexlmax)
func (asnInfoService *AsnInfoService) GetAsnInfoInfoList(info cfscanReq.AsnInfoSearch) (list []cfscan.AsnInfo, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&cfscan.AsnInfo{})
	var asnInfos []cfscan.AsnInfo
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	//模糊查询
	if info.AsnName != "" {
		db = db.Where("asn_name LIKE ?", "%"+info.AsnName+"%")
	}
	if info.FullName != "" {
		db = db.Where("full_name LIKE ?", "%"+info.FullName+"%")
	}
	if info.AllocationCountry != "" {
		db = db.Where("allocation_country LIKE ?", "%"+info.AllocationCountry+"%")

	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&asnInfos).Error
	return asnInfos, total, err
}
