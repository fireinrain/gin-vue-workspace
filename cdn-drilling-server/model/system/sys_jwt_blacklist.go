package system

import (
	"github.com/fireinrain/gin-vue-workspace/cdn-drilling-server/global"
)

type JwtBlacklist struct {
	global.GVA_MODEL
	Jwt string `gorm:"type:text;comment:jwt"`
}
