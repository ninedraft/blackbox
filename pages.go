package main

import (
	"errors"
	"html/template"
	"io"
	"strings"
	"sync"
)

// pages map must be changed to LRU cache

var (
	pages         = map[string]Pager{}
	pagesMutex    = &sync.Mutex{}
	ErrNoSuchPage = errors.New("now such page")
)

func RoutePage(id string, wr io.Writer) error {
	if pages[id] == nil {
		return ErrNoSuchPage
	}
	page := pages[id]
	loweredID := strings.ToLower(id)
	switch {
	case loweredID == "index" || loweredID == "home":
		if Templates["index"] != nil {
			return RenderPage(FillPage(pages["index"]), Templates["index"], wr)
		}
		fallthrough
	default:
		return RenderPage(FillPage(page, PageDate.ToPage()), Templates["article"], wr)
	}
}

func RenderPage(pd PageData, templ *template.Template, wr io.Writer) error {
	return templ.Execute(wr, pd)
}

func FillPage(page Pager, data ...PageData) PageData {
	return PageData{}.Merge(data...).Merge(page.ToPage())
}

func SetPage(id string, page Pager) {
	pagesMutex.Lock()
	defer pagesMutex.Unlock()
	pages[id] = page
}
