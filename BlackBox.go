// BlackBox project BlackBox.go
package main

import (
	"errors"
	"flag"
	"html/template"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	ErrIsNotDir           = errors.New("is not dir")
	ErrInvalidArticleFile = errors.New("invalid article file")
)

type BlackBox struct {
	articles     map[string]*Article
	articlepath  string
	templatepath string
	content      map[string]io.WriterTo
	templates    map[string]*template.Template
	logger       Logger
	mux          *mux.Router
}

func (bb *BlackBox) ParseTemplates() error {
	templfiles, err := GetFilesIfDir(bb.templatepath)
	if err != nil {
		return err
	}
	var t *template.Template
	for _, tfile := range templfiles {
		bb.LogPrintf("Loading template %q:...", tfile.Name())
		t, err = template.ParseFiles(bb.templatepath + "/" + tfile.Name())
		if err != nil {
			bb.LogPrintf("FAILED!:\n%v\n", err)
		} else {
			bb.templates[FileNameNoExtension(tfile.Name())] = t
			bb.LogPrintf("OK\n")
		}
	}
	return err
}

func (bb *BlackBox) LoadArticles() error {
	finfos, err := GetFilesIfDir(bb.articlepath)
	if err != nil {
		return err
	}
	var article *Article
	var articletempl *template.Template
	for _, f := range finfos {
		article, err = NewArticleFromFile(bb.articlepath + "/" + f.Name())
		bb.LogPrintf("Try to load article %q...", f.Name())
		if err != nil {
			bb.LogPrintf("FAILED!\n%v\n", err)
			continue
		}
		err = article.Reload()
		if err != nil {
			bb.LogPrintf("FAILED! Error while reloading\n%v\n", err)
			continue
		}
		articletempl = bb.templates[article.Template]
		if articletempl == nil {
			bb.LogPrintf("FAILED!\nCan't find template %q\n", article.Template)
			continue
		}
		article.SetTemplate(articletempl)
		err = article.Render()
		if err != nil {
			bb.LogPrintf("FAILED! Error while rendering\n%v\n", err)
		}
		bb.LogPrintf("OK:\n\tID: %q\n\tTitle: %q\n\tTags: %v\n\tTemplate: %q\n", article.ID, article.Title, article.Tags, article.Template)
		bb.articles[article.ID] = article
	}
	return nil
}

func (blog *BlackBox) SetArticleHandler() {
	blog.mux.HandleFunc("/articles/{articleID}", func(resp http.ResponseWriter, req *http.Request) {
		articleID := mux.Vars(req)["articleID"]
		blog.LogPrintf("New request: %q...", articleID)
		article := blog.articles[articleID]
		if article == nil {
			blog.LogPrintf("NOT FOUND %q\n", articleID)
			resp.WriteHeader(http.StatusNotFound)
			resp.Write([]byte("NOT FOUND ☹\n☄"))
			return
		}
		n, err := article.WriteTo(resp)
		if err != nil {
			blog.LogPrintln(err)
			return
		}
		blog.LogPrintf("Ok: %d bytes\n", n)
	})

}

func (bb *BlackBox) LogPrintf(f string, v ...interface{}) {
	bb.logger.Printf(f, v...)
}

func (bb *BlackBox) LogPrintln(v ...interface{}) {
	bb.logger.Println(v...)
}

func (bb *BlackBox) Mux() *mux.Router {
	return bb.mux
}

type BlackBoxConfig struct {
	ArticlePath  string
	TemplatePath string
	StaticPath   string
	Logger       Logger
}

func NewBlackBox(config BlackBoxConfig) *BlackBox {
	nbb := &BlackBox{
		articles:     map[string]*Article{},
		articlepath:  config.ArticlePath,
		templatepath: config.TemplatePath,
		logger:       config.Logger,
		mux:          mux.NewRouter(),
		templates:    map[string]*template.Template{},
	}
	nbb.SetArticleHandler()
	fileserver := http.StripPrefix("/static/", http.FileServer(http.Dir(config.StaticPath)))
	nbb.mux.PathPrefix("/static/").Handler(fileserver)
	return nbb
}

func main() {
	var articlepath string
	var templatepath string
	var staticpath string
	flag.StringVar(&articlepath, "a", "./articles", `path to markdown documents with articles`)
	flag.StringVar(&templatepath, "t", "./templates", `path to html templates`)
	flag.StringVar(&staticpath, "s", "./static", `path to static content`)
	flag.Parse()
	blog := NewBlackBox(BlackBoxConfig{
		ArticlePath:  articlepath,
		TemplatePath: templatepath,
		Logger:       &DebugLogger{},
		StaticPath:   staticpath,
	})
	err := blog.ParseTemplates()
	if err != nil {
		blog.LogPrintf("%v\n", err)
		return
	}
	err = blog.LoadArticles()
	if err != nil {
		blog.LogPrintf("%v\n", err)
		return
	}
	http.Handle("/", blog.Mux())
	blog.LogPrintln(http.ListenAndServe(":8080", nil))
}
