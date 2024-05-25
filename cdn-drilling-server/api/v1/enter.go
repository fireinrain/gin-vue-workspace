package v1

import (
	"github.com/fireinrain/gin-vue-workspace/cdn-drilling-server/api/v1/example"
	"github.com/fireinrain/gin-vue-workspace/cdn-drilling-server/api/v1/system"
)

type ApiGroup struct {
	SystemApiGroup  system.ApiGroup
	ExampleApiGroup example.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
