package request

import (
	"github.com/fireinrain/gin-vue-workspace/cdn-drilling-server/model/common/request"
	"github.com/fireinrain/gin-vue-workspace/cdn-drilling-server/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
