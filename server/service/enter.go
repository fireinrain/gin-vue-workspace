package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/cfscan"
	"github.com/flipped-aurora/gin-vue-admin/server/service/example"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
	CfscanServiceGroup  cfscan.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
