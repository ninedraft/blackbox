// page.go
package main

import (
	"bytes"
	"html/template"
	"io"
	"strings"
	"time"

	"github.com/golang-commonmark/markdown"
	yaml "gopkg.in/yaml.v2"
)

type Page struct {
	Title   string   `yaml: id`
	Tags    []string `yaml: tags`
	Body    template.HTML
	Created time.Time
	Author  string `yaml: author`
}

func (page *Page) ReadFrom(r io.Reader) (int64, error) {
	buf := &bytes.Buffer{}
	n, err := buf.ReadFrom(r)
	if err != nil && err != io.EOF {
		return 0, err
	}
	raw := buf.String()
	data := strings.SplitN(raw, "---", 3)
	if len(data) < 2 {
		return 0, ErrInvalidArticleFile
	}
	err = yaml.Unmarshal([]byte(data[1]), page)
	if err != nil {
		return 0, err
	}
	body := data[2]
	md := markdown.New(markdown.HTML(true))
	page.Body = template.HTML(md.RenderToString([]byte(body)))
	return n, nil
}
