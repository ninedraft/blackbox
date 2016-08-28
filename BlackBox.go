// BlackBox project BlackBox.go
package main

// filesystem:
// ~
// ----/templates
// ----/cache
// ----/articles
// ----/media
// ----/apps

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var (
	Templates        = map[string]*template.Template{}
	Log       Logger = new(LoggerDebug)
)

func IndexHandler(rw http.ResponseWriter, r *http.Request) {
	_, err := rw.Write([]byte("Hello, it's index page!"))
	if err != nil {
		log.Printf("error while writing response on index:  %v\n", err)
	}
}

func main() {
	Port := flag.Uint64("port", 8080, "port")
	TemplatePath := flag.String("template", "./data/templates/", "path to templates")
	StaticPath := flag.String("media", "./data/static/", "path to media content")
	ArticlesPath := flag.String("articles", "./data/articles/", "path to articles")
	AppsPath := flag.String("apps", "./data/apps/", "path to apps")
	flag.Parse()
	Log.Printf("templates: %s\nstatic: %s\narticles: %s\napps: %s\nserving at port %d\n", *TemplatePath, *StaticPath, *ArticlesPath, *AppsPath, *Port)

	router := mux.NewRouter()

	absStaticPath, err := filepath.Abs(*StaticPath)
	Ok("error while calculating absolute static path: %v\n", err)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(absStaticPath))))

	indexArticleFile, err := os.Open(*ArticlesPath + "index.md")
	Ok("error while loading index article: %v\n", err)
	indexArticle, err := MarkdownArticleFromReader(indexArticleFile)
	Ok("error while reading index article: %v\n", err)
	Ok("error while closing index article file: %v\n", indexArticleFile.Close())
	pages["index"] = indexArticle

	indexTemplate, err := template.ParseFiles(*TemplatePath + "index.html")
	Ok("error while loading index template: %v\n", err)
	Templates["index"] = indexTemplate

	router.HandleFunc("/article/{pageid}", PageHandler)
	router.HandleFunc("/index", PageHandleFunc("index"))
	router.HandleFunc("/home", PageHandleFunc("index"))

	server := http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:" + strconv.FormatUint(*Port, 10),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println(server.ListenAndServe())
}
