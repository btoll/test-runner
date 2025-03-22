{{- range . -}}
{{- if eq .Succeeded false -}}
Suite: {{ .Description }}
Path: {{ .Path }}
Succeeded: {{ .Succeeded }}
PreRunStats:
    TotalSpecs: {{ .PreRunStats.TotalSpecs }}
    SpecsThatWillRun: {{ .PreRunStats.SpecsThatWillRun -}}
{{- end }}
{{- range .SpecReports -}}
{{- if eq .State "failed" }}
Specs:
    - Name: {{ .Name }}
      Type: {{ .Type }}
      State: {{ .State }}
      Attempts: {{- .Attempts }}
          LineNumber: {{ .Failure.Location.LineNumber }}
          StackTrace: {{ .Failure.Location.StackTrace -}}
{{- end -}}
{{- end -}}
{{- end }}

