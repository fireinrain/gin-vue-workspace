package response

import "github.com/fireinrain/gin-vue-workspace/cdn-drilling-server/model/example"

type ExaFileResponse struct {
	File example.ExaFileUploadAndDownload `json:"file"`
}
