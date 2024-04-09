---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .LoggerName }}
  namespace: {{ .HelmDeployNamespace }}
  labels:
    app: {{ .AppLabel }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{ .AppLabel }}
  template:
    metadata:
      labels:
        app: {{ .AppLabel }}
    spec:
      containers:
      - name: logger
        # A custom agnhost image (eu.gcr.io/gardener-project/3rd/agnhost) is used instead of the upstream one (registry.k8s.io/e2e-test-images/agnhost)
        # because this Deployment is created in a Seed cluster and the image needs to be signed with particular keys.
        image: eu.gcr.io/gardener-project/3rd/agnhost:2.40
        command: ["/bin/sh"]
        args:
          - -c
          - |-
{{ if .DeltaLogsCount }}
            /agnhost logs-generator --log-lines-total={{ .DeltaLogsCount }} --run-duration={{ .DeltaLogsDuration }}
{{- end }}
            /agnhost logs-generator --log-lines-total={{ .LogsCount }} --run-duration={{ .LogsDuration }}

            # Sleep forever to prevent restarts
            while true; do
              sleep 3600;
            done
        resources:
          limits:
            cpu: 8m
            memory: 30Mi
          requests:
            cpu: 4m
            memory: 10Mi
      securityContext:
        fsGroup: 65532
        runAsUser: 65532
        runAsNonRoot: true
{{ if .NodeName }}
      nodeName: {{ .NodeName }}
{{- end }}
{{ if .NodeSelector }}
      nodeSelector: {{ .NodeSelector }}
{{- end }}
