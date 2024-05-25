package service

import (
	"github.com/fireinrain/gin-vue-workspace/cdn-drilling-server/service/example"
	"github.com/fireinrain/gin-vue-workspace/cdn-drilling-server/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
