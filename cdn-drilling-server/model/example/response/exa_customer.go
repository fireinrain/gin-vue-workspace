package response

import "github.com/fireinrain/gin-vue-workspace/cdn-drilling-server/model/example"

type ExaCustomerResponse struct {
	Customer example.ExaCustomer `json:"customer"`
}
