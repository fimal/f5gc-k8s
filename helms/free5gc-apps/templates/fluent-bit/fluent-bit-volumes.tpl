{{/*
volumes configuration for pods where
 one container (enforcer) writes logs
 and another (fluent-bit) reads them
*/}}
{{- define "fluent-bit-volumes" }}
- name: securebeat-ca-volume-{{ .Values.global.app.name }}
  configMap:
    name: {{ .Values.global.namesPrefix }}securebeat-ca-config-{{ .Values.global.app.name }}
- name: securebeat-client-volume-{{ .Values.global.app.name }}
  secret:
    secretName: {{ .Values.global.namesPrefix }}securebeat-client-secret-{{ .Values.global.app.name }}
- name: logs-volume
  emptyDir: {}
{{- end -}}