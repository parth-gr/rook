{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Define imagePullSecrets option to pass to all service accounts
*/}}
{{- define "imagePullSecrets" }}
{{- if .Values.imagePullSecrets -}}
imagePullSecrets:
{{ toYaml .Values.imagePullSecrets }}
{{- end -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "rook-ceph.labels" -}}
app.kubernetes.io/name: rook-ceph
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: csi
app.kubernetes.io/part-of: rook-ceph
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/created-by: rook-ceph-operator
{{- end -}}
