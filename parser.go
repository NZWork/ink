package main

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func ParseMarkdownFileToHTML(md []byte) []byte {
	unsafe := blackfriday.MarkdownCommon(md)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return html
}
