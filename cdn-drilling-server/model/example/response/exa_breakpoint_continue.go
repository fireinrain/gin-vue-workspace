package response

import "github.com/fireinrain/gin-vue-workspace/cdn-drilling-server/model/example"

type FilePathResponse struct {
	FilePath string `json:"filePath"`
}

type FileResponse struct {
	File example.ExaFile `json:"file"`
}
