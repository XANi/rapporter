package web

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
	"html"
)

var bluePolicy = bluemonday.StrictPolicy()

func (b *WebBackend) markdownParse(s string) string {
	md := bluePolicy.SanitizeBytes([]byte(s))
	// turn back lt and gt into proper characters
	md = []byte(html.UnescapeString(string(md)))
	extensions := parser.CommonExtensions
	p := parser.NewWithExtensions(extensions)
	mdhtml := markdown.ToHTML(md, p, nil)
	return string(mdhtml)
}
