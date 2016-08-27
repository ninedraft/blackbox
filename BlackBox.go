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
	TemplatePath := flag.String("template", "./templates", "path to templates")
	StaticPath := flag.String("media", "./static", "path to media content")
	ArticlesPath := flag.String("articles", "./articles", "path to articles")
	AppsPath := flag.String("apps", "./apps", "path to apps")
	flag.Parse()
	Log.Printf("templates: %s\nstatic: %s\narticles: %s\napps: %s\nserving at port %d\n", *TemplatePath, *StaticPath, *ArticlesPath, *AppsPath, *Port)
	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
	router.HandleFunc("/index", IndexHandler)
	server := http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:" + strconv.FormatUint(*Port, 10),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println(server.ListenAndServe())
}
