package web

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
)

func (b *WebBackend) markdownParse(s string) string {
	md := bluemonday.StrictPolicy().SanitizeBytes([]byte(s))
	extensions := parser.CommonExtensions
	p := parser.NewWithExtensions(extensions)
	html := markdown.ToHTML(md, p, nil)
	return string(html)
}
