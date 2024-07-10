package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/cfscan"
	"gorm.io/gorm"
)

func bizModel(db *gorm.DB) error {
	return db.AutoMigrate(cfscan.AsnInfo{}, cfscan.SubmitScan{}, cfscan.ScheduleTask{}, cfscan.ScheduleTaskHist{})
}
