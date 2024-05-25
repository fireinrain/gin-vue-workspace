package response

import "github.com/fireinrain/gin-vue-workspace/cdn-drilling-server/model/system"

type SysAPIResponse struct {
	Api system.SysApi `json:"api"`
}

type SysAPIListResponse struct {
	Apis []system.SysApi `json:"apis"`
}
