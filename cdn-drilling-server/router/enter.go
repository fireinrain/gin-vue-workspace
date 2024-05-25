package router

import (
	"github.com/fireinrain/gin-vue-workspace/cdn-drilling-server/router/example"
	"github.com/fireinrain/gin-vue-workspace/cdn-drilling-server/router/system"
)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
