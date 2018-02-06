package main

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	adsabsurl    = `http://adsabs.harvard.edu/cgi-bin/nph-abs_connect`
	adsbiburl    = `http://adsabs.harvard.edu/cgi-bin/nph-bib_query`
	validlinks   = `ACDEFGNORSTUX`
	abstagbefore = `Abstract</h3>`
	abstagafter  = `<hr/>`
	pattern      = `(?m)` + abstagbefore + `[\s\S]*?` + abstagafter
)

func GetLinks(form *Form) ([]*Paper, error) {
	doc, err := SendForm(form)
	if err != nil {
		return nil, err
	}
	return GetLinksFromDocument(doc)
}

func GetAbstract(paper *Paper) (string, error) {
	if paper.HasLink(`A`) {
		res, err := http.Get(paper.urls[`A`])
		if err != nil {
			return "", err
		}
		defer res.Body.Close()
		doc, err := goquery.NewDocumentFromResponse(res)
		if err != nil {
			return "", err
		}
		r := regexp.MustCompile(pattern)
		s, err := doc.Find("body").Html()
		s = r.FindString(s)
		// Trim unnecessary chars before&after abstract
		s = strings.Split(s, abstagbefore)[1]
		s = s[0 : len(s)-5]
		return s, nil
	}
	return "", nil
}

//func GetAbstractFromBibcode(bibcode string) (string, error) {
//	res, err := http.PostForm()
//}

func GetBibTexEntry(bibcode string) (string, error) {
	values := url.Values{}
	values.Add(`bibcode`, bibcode)
	values.Add(`data_type`, `BIBTEX`)
	values.Add(`db_key`, `AST`)
	values.Add(`nocookieset`, `1`)
	res, err := http.PostForm(adsbiburl, values)
	defer res.Body.Close()
	if err != nil {
		return "", err
	}
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return "", err
	}
	bibtex := `@` + strings.TrimSpace(strings.Split(doc.Find("body").Text(), `@`)[1])
	return bibtex, nil
}

func SendForm(form *Form) (*goquery.Document, error) {
	res, err := http.PostForm(adsabsurl, form.values)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromResponse(res)
}

func GetLinksFromDocument(doc *goquery.Document) ([]*Paper, error) {
	// Get bibcodes
	bibcodes_str, _ := doc.Find("form > input").Attr("value")
	bibcodes := strings.Split(bibcodes_str, ";")
	// Get links
	table := doc.Find("form > table").Eq(1)
	papers := make([]*Paper, len(bibcodes), len(bibcodes))
	cnt := 0
	table.Find("tbody>tr").Each(func(i int, s *goquery.Selection) {
		// Links
		if i > 0 && i%3 == 0 {
			papers[cnt] = NewPaper()
			linktypes := s.Find("td").Last().Find("a")
			papers[cnt].bibcode = bibcodes[cnt]
			linktypes.Each(func(_ int, s *goquery.Selection) {
				link, _ := s.Attr("href")
				papers[cnt].SetURL(link, s.Text())
			})
		}
		// Authors and Titles
		if i > 1 && i%3 == 1 {
			td := s.Find("td")
			authors := td.Eq(1).Text()
			title := td.Eq(3).Text()
			papers[cnt].authors = authors
			papers[cnt].title = title
			cnt++
		}
	})
	return papers, nil
}
