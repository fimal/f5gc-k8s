{{/*
fluentbit container manifest
reads enforcer logs (security + access) and sends over http(s) to logstash
*/}}
{{ define "fluent-bit-container" }}
- name: fluentbit
  image: {{ include "image.registry" . | coalesce .Values.containers.fluentbit.registry }}/waas-fluentbit{{ include "image.version" . | coalesce .Values.containers.fluentbit.version }}
  imagePullPolicy: {{ .Values.global.image.pullPolicy }}
  env:
  - name: FB_LS_HOST
    value : {{ .Values.global.namesPrefix }}logstash-service.{{ required "A valid KWAF Namespace required!" .Values.global.kwafns  }}.svc.cluster.local
  - name: FB_PREFIX
    value: "{{ .Values.global.namesPrefix }}{{ .Chart.Name }}-{{ .Values.global.app.name}}.{{ .Release.Namespace }}"
  resources:
    {{- toYaml .Values.containers.fluentbit.resources | nindent 4 }}
  volumeMounts:
  - name: securebeat-client-volume-{{ .Values.global.app.name }}
    mountPath: /etc/securebeat/client
  - name: securebeat-ca-volume-{{ .Values.global.app.name }}
    mountPath: /etc/securebeat/ca
  - name: logs-volume
    mountPath: /logs
{{- end -}}
