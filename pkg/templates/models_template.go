package templates

const modelsTemplate = `package {{.PackageName}}
{{if .NeedsTime}}
import (
	"time"
)
{{end}}

{{range .TypeAliases}}
{{if .Description}}// {{.Name}} {{.Description}}{{end}}
type {{.Name}} {{.Type}}
{{end}}

{{range .Models}}
{{if .Description}}// {{.Name}} {{.Description}}{{end}}
type {{.Name}} struct {
{{range .Fields}}	{{.Name}} {{.Type}} ` + "`json:\"{{.JSONTag}}\"`" + `{{if .Description}} // {{.Description}}{{end}}
{{end}}}
{{end}}
`