package global

{{- if .HasGlobal }}

import "github.com/fireinrain/gin-vue-workspace/cdn-drilling-server/plugin/{{ .Snake}}/config"

var GlobalConfig = new(config.{{ .PlugName}})
{{ end -}}