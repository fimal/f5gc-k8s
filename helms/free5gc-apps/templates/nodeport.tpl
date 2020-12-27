{{/*
Include a service node port
*/}}
{{- define "waas.nodeport" }}
{{- if eq .Values.global.serviceType "NodePort" }}
{{- $port := toString .Values.global.nodePort | int }}
{{- if  $port }}
{{- printf "nodePort: %d" $port }}
{{- end }}
{{- end }}
{{- end }}