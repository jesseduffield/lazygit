package flaggy

// defaultHelpTemplate is the help template used by default
// {{if (or (or (gt (len .StringFlags) 0) (gt (len .IntFlags) 0)) (gt (len .BoolFlags) 0))}}
// {{if (or (gt (len .StringFlags) 0) (gt (len .BoolFlags) 0))}}
const defaultHelpTemplate = `{{.CommandName}}{{if .Description}} - {{.Description}}{{end}}{{if .PrependMessage}}
{{.PrependMessage}}{{end}}
{{if .UsageString}}
  Usage:
    {{.UsageString}}{{end}}{{if .Positionals}}

  Positional Variables: {{range .Positionals}}
    {{.Name}}  {{.Spacer}}{{if .Description}} {{.Description}}{{end}}{{if .DefaultValue}} (default: {{.DefaultValue}}){{else}}{{if .Required}} (Required){{end}}{{end}}{{end}}{{end}}{{if .Subcommands}}

  Subcommands: {{range .Subcommands}}
    {{.LongName}}{{if .ShortName}} ({{.ShortName}}){{end}}{{if .Position}}{{if gt .Position 1}}  (position {{.Position}}){{end}}{{end}}{{if .Description}}   {{.Spacer}}{{.Description}}{{end}}{{end}}
{{end}}{{if (gt (len .Flags) 0)}}
  Flags: {{if .Flags}}{{range .Flags}}
    {{if .ShortName}}-{{.ShortName}} {{else}}   {{end}}{{if .LongName}}--{{.LongName}}{{end}}{{if .Description}}   {{.Spacer}}{{.Description}}{{if .DefaultValue}} (default: {{.DefaultValue}}){{end}}{{end}}{{end}}{{end}}
{{end}}{{if .AppendMessage}}{{.AppendMessage}}
{{end}}{{if .Message}}
{{.Message}}{{end}}
`
