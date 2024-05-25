package system

import (
	"github.com/fireinrain/gin-vue-workspace/cdn-drilling-server/config"
)

// 配置文件结构体
type System struct {
	Config config.Server `json:"config"`
}
