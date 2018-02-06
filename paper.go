package main

import (
	"fmt"
	"strings"
)

/*=======================================================
/*                    Paper
/*=======================================================*/
type Paper struct {
	bibcode  string
	title    string
	authors  string
	abstract string
	urls     map[string]string
}

func NewPaper() *Paper {
	return &Paper{
		bibcode: "",
		urls:    make(map[string]string)}
}

func (paper *Paper) Show() {
	for key, val := range paper.urls {
		if val != "" {
			fmt.Println(key + ":" + val)
		}
	}
}

func (paper *Paper) GetURL(linktype string) string {
	if paper.HasLink(linktype) {
		return paper.urls[linktype]
	}
	return ""
}

func (paper *Paper) SetURL(url string, linktype string) error {
	if len(linktype) == 1 && strings.Contains(validlinks, linktype) {
		paper.urls[strings.ToUpper(linktype)] = url
		return nil
	}
	err := fmt.Errorf(`linktype must be a single character in %v`, validlinks)
	return err
}

func (paper *Paper) AvailableLinkTypes() string {
	linktypes := ""
	for _, linktype := range validlinks {
		s := string(linktype)
		if paper.HasLink(s) {
			linktypes += s
		}
	}
	return linktypes
}

func (paper *Paper) AvailableLinkTypesIn(linktypes string) string {
	availablelinks := ""
	for _, linktype := range validlinks {
		s := string(linktype)
		if paper.HasLink(s) && strings.Contains(linktypes, s) {
			availablelinks += s
		}
	}
	return availablelinks
}

func (paper *Paper) HasLinkAny(linktypes string) bool {
	for _, linktype := range linktypes {
		if paper.HasLink(string(linktype)) {
			return true
		}
	}
	return false
}

func (paper *Paper) HasLink(linktype string) bool {
	linktype = strings.ToUpper(linktype)
	if len(linktype) == 1 && strings.Contains(validlinks, linktype) {
		if paper.urls[linktype] == "" {
			return false
		} else {
			return true
		}
	}
	fmt.Println(`linktype must be a single character in %v`, validlinks)
	return false
}
