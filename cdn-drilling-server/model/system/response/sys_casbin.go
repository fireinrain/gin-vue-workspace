package response

import (
	"github.com/fireinrain/gin-vue-workspace/cdn-drilling-server/model/system/request"
)

type PolicyPathResponse struct {
	Paths []request.CasbinInfo `json:"paths"`
}
