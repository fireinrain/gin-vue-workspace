package cfscan

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cfscan"
	cfscanReq "github.com/flipped-aurora/gin-vue-admin/server/model/cfscan/request"
)

type SubmitScanService struct{}

// CreateSubmitScan 创建submitScan表记录
// Author [piexlmax](https://github.com/piexlmax)
func (submitScanService *SubmitScanService) CreateSubmitScan(submitScan *cfscan.SubmitScan) (err error) {
	err = global.GVA_DB.Create(submitScan).Error
	return err
}

// DeleteSubmitScan 删除submitScan表记录
// Author [piexlmax](https://github.com/piexlmax)
func (submitScanService *SubmitScanService) DeleteSubmitScan(ID string) (err error) {
	err = global.GVA_DB.Delete(&cfscan.SubmitScan{}, "id = ?", ID).Error
	return err
}

// DeleteSubmitScanByIds 批量删除submitScan表记录
// Author [piexlmax](https://github.com/piexlmax)
func (submitScanService *SubmitScanService) DeleteSubmitScanByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]cfscan.SubmitScan{}, "id in ?", IDs).Error
	return err
}

// UpdateSubmitScan 更新submitScan表记录
// Author [piexlmax](https://github.com/piexlmax)
func (submitScanService *SubmitScanService) UpdateSubmitScan(submitScan cfscan.SubmitScan) (err error) {
	err = global.GVA_DB.Model(&cfscan.SubmitScan{}).Where("id = ?", submitScan.ID).Updates(&submitScan).Error
	return err
}

// GetSubmitScan 根据ID获取submitScan表记录
// Author [piexlmax](https://github.com/piexlmax)
func (submitScanService *SubmitScanService) GetSubmitScan(ID string) (submitScan cfscan.SubmitScan, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&submitScan).Error
	return
}

// GetSubmitScanInfoList 分页获取submitScan表记录
// Author [piexlmax](https://github.com/piexlmax)
func (submitScanService *SubmitScanService) GetSubmitScanInfoList(info cfscanReq.SubmitScanSearch) (list []cfscan.SubmitScan, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&cfscan.SubmitScan{})
	var submitScans []cfscan.SubmitScan
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.ScanDesc != "" {
		db = db.Where("scan_desc LIKE ?", "%"+info.ScanDesc+"%")
	}
	if info.ScanType != "" {
		db = db.Where("scan_type = ?", info.ScanType)
	}
	if info.AsnNumber != "" {
		db = db.Where("asn_number LIKE ?", "%"+info.AsnNumber+"%")
	}
	if info.ScanStatus != "" {
		db = db.Where("scan_status = ?", info.ScanStatus)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&submitScans).Error
	return submitScans, total, err
}
