package response

import "github.com/fireinrain/gin-vue-workspace/cdn-drilling-server/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
