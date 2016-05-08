// MDArticle.go
package main

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"strings"

	"github.com/golang-commonmark/markdown"
	yaml "gopkg.in/yaml.v2"
)

var (
	ErrTemplateNotFound = errors.New("template not found")
)

type Article struct {
	ID       string   `yaml: id`
	Title    string   `yaml: title`
	Author   string   `yaml: author`
	Tags     []string `yaml: tags`
	Template string   `yaml: template`
	Body     template.HTML
	rawbody  string
	Preview  template.HTML

	path              string
	compilledtemplate *template.Template
	buf               []byte
}

func NewArticleFromFile(artpath string) (*Article, error) {
	narticle := &Article{
		path: artpath,
	}
	return narticle, nil
}

func (article *Article) SetTemplate(t *template.Template) {
	article.compilledtemplate = t
}

func (article *Article) Render() error {
	if article.compilledtemplate == nil {
		return ErrTemplateNotFound
	}
	md := markdown.New(markdown.HTML(true))
	article.Body = template.HTML(md.RenderToString([]byte(article.rawbody)))
	buf := bytes.NewBuffer(article.buf)
	err := article.compilledtemplate.Execute(buf, article)
	if err != nil {
		return err
	}
	if len(article.buf) < buf.Len() {
		article.buf = append(article.buf, make([]byte, 2*buf.Len())...)
	}
	copy(article.buf, buf.Bytes())
	return nil
}

func (mda *Article) ParseFile(name string) error {
	bin, err := ioutil.ReadFile(name)
	if err == nil {
		data := strings.SplitN(string(bin), "---", 3)
		if len(data) < 2 {
			return ErrInvalidArticleFile
		}
		err = yaml.Unmarshal([]byte(data[1]), mda)
		if err != nil {
			return err
		}
		mda.rawbody = data[2]
	}
	return err
}

func (article *Article) Reload() error {
	err := article.ParseFile(article.path)
	return err
}

func (mda *Article) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(mda.buf)
	if err == io.EOF {
		return int64(n), nil
	}
	return int64(n), err
}
