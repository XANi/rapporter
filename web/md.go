package web

import (
	"github.com/gomarkdown/markdown"
	mdhtml "github.com/gomarkdown/markdown/html"
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
	// parser can't be cached, must be one instance per ToHTML call... for some reason
	p := parser.NewWithExtensions(extensions)
	htmlRenderer := mdhtml.NewRenderer(mdhtml.RendererOptions{})
	mdhtml := markdown.ToHTML(md, p, htmlRenderer)
	return string(mdhtml)
}
