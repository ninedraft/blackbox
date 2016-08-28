// page.go
package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"regexp"
	"strings"

	"github.com/golang-commonmark/markdown"
)

const (
	TypeMarkdownArticle = PageType("markdown article")
)

type MarkdownArticle struct {
	ID    string
	Title string
	Body  template.HTML
	Tags  []string
}

func (markdownArticle *MarkdownArticle) ToPage() PageData {
	pd := PageData{}
	pd["Title"] = markdownArticle.Title
	pd["Tags"] = markdownArticle.Tags
	pd["ID"] = markdownArticle.ID
	pd["Body"] = markdownArticle.Body
	return pd
}

func (markdownArticle *MarkdownArticle) Type() PageType {
	return TypeMarkdownArticle
}

func MarkdownArticleFromReader(r io.Reader) (*MarkdownArticle, error) {
	buf := &bytes.Buffer{}
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("error while reading page: %v", err)
	}
	markdownArticle := &MarkdownArticle{}
	idLine := regexp.MustCompile("^[\\s, \\t, \\n, \\r]*(\\w+)").Find(buf.Bytes())
	markdownArticle.ID = strings.TrimSpace(string(idLine))
	titleLine := regexp.MustCompile("#[\\s, \\t]+([^\\n]+)").Find(buf.Bytes())
	title := strings.TrimLeftFunc(string(titleLine), func(r rune) bool {
		return r == '#' || r == ' ' || r == '	'
	})
	markdownArticle.Title = strings.TrimSpace(title)
	md := markdown.New(markdown.HTML(true))
	body := md.RenderToString(buf.Bytes())
	markdownArticle.Body = template.HTML(body)
	markdownArticle.Tags = ExctractTags(body)
	return markdownArticle, nil
}
