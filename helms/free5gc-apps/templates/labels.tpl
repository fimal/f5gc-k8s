{{/*
Include global labels
*/}}
{{- define "waas.labels" }}
{{- range $k, $v := .Values.global.labels }}
{{ $k }} : {{ $v}}
{{- end }}
{{- end }}
