package template

import _ "embed"

//go:embed template.html
var t1Html string

func GetHtmlTemplate() string {
	return t1Html
}
