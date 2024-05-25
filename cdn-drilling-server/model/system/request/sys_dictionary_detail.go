package request

import (
	"github.com/fireinrain/gin-vue-workspace/cdn-drilling-server/model/common/request"
	"github.com/fireinrain/gin-vue-workspace/cdn-drilling-server/model/system"
)

type SysDictionaryDetailSearch struct {
	system.SysDictionaryDetail
	request.PageInfo
}
