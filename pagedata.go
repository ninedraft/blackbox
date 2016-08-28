package main

import (
	"time"
)

var (
	PageDate = pageDate{}
)

type PageType string

const (
	TypeIndex = PageType("index")
	TypeDate  = PageType("date")
)

type PageData map[string]interface{}

func (pageData PageData) Copy() PageData {
	pd := PageData{}
	for i, v := range pageData {
		pd[i] = v
	}
	return pd
}

func (pageData PageData) Merge(pd ...PageData) PageData {
	res := pageData.Copy()
	for _, d := range pd {
		for i, v := range d {
			res[i] = v
		}
	}
	return res
}

type Pager interface {
	ToPage() PageData
	Type() PageType
}

type pageDate struct{}

func (pdate *pageDate) ToPage() PageData {
	return PageData{
		"Date": time.Now().Format("2006 Jan Mon 15:04"),
	}
}

func (pdate *pageDate) Type() PageType {
	return TypeDate
}
