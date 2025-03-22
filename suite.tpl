{{ range . -}}
Suite: {{.Description}}
Path: {{.Path}}
Succeeded: {{.Succeeded}}
PreRunStats:
    TotalSpecs: {{.PreRunStats.TotalSpecs}}
    SpecsThatWillRun: {{.PreRunStats.SpecsThatWillRun}}
Specs:
{{- range .SpecReports }}
    - Name: {{.Name}}
      Type: {{.Type}}
      State: {{.State}}
      Attempts: {{.Attempts}}
      {{- if eq .State "failed" }}
          LineNumber: {{.Failure.Location.LineNumber}}
          StackTrace: {{.Failure.Location.StackTrace}}
      {{- end }}
{{ end }}
{{- end }}

