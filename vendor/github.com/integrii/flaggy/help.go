package flaggy

// defaultHelpTemplate is the help template used by default
// {{if (or (or (gt (len .StringFlags) 0) (gt (len .IntFlags) 0)) (gt (len .BoolFlags) 0))}}
// {{if (or (gt (len .StringFlags) 0) (gt (len .BoolFlags) 0))}}
const defaultHelpTemplate = `{{range $idx, $line := .Lines}}{{if gt $idx 0}}
{{end}}{{$line}}{{end}}`
