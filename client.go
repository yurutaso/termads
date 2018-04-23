package termads

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	ADS_ABS_URL    = `http://adsabs.harvard.edu/cgi-bin/nph-abs_connect`
	ADS_BIB_URL    = `http://adsabs.harvard.edu/cgi-bin/nph-bib_query`
	ABSTAG_BEFORE  = `Abstract</h3>`
	ABSTAG_AFTER   = `<hr/>`
	SEARCH_PATTERN = `(?m)` + ABSTAG_BEFORE + `[\s\S]*?` + ABSTAG_AFTER
)

func GetPapers(form *Form) ([]Paper, error) {
	fmt.Printf("Waiting for response from ADS.\n")
	res, err := http.PostForm(ADS_ABS_URL, form.values)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	fmt.Printf("Extracting papers from response.\n")
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}
	return GetPapersFromDocument(doc)
}

func GetPapersFromDocument(doc *goquery.Document) ([]Paper, error) {
	// Get bibcodes
	bibcodes_str, _ := doc.Find("form > input").Attr("value")
	bibcodes := strings.Split(bibcodes_str, ";")

	// Get links for each bibcode
	table := doc.Find("form > table").Eq(1)
	papers := make([]Paper, len(bibcodes), len(bibcodes))
	cnt := 0
	table.Find("tbody>tr").Each(func(i int, s *goquery.Selection) {
		// Links
		if i > 0 && i%3 == 0 {
			papers[cnt] = NewPaper()
			linktypes := s.Find("td").Last().Find("a")
			papers[cnt].SetBibcode(bibcodes[cnt])
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
			papers[cnt].SetAuthors(authors)
			papers[cnt].SetTitle(title)
			cnt++
		}
	})
	return papers, nil
}

func GetAbstract(_url string) (string, error) {
	res, err := http.Get(_url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return "", err
	}
	r := regexp.MustCompile(SEARCH_PATTERN)
	s, err := doc.Find("body").Html()
	s = r.FindString(s)
	// Trim unnecessary chars before&after abstract
	s = strings.Split(s, ABSTAG_BEFORE)[1]
	s = s[0 : len(s)-5]
	return s, nil
}

func GetBibTex(bibcode string) (string, error) {
	values := url.Values{}
	values.Add(`bibcode`, bibcode)
	values.Add(`data_type`, `BIBTEX`)
	values.Add(`db_key`, `AST`)
	values.Add(`nocookieset`, `1`)
	res, err := http.PostForm(ADS_BIB_URL, values)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return "", err
	}

	bibtex := `@` + strings.TrimSpace(strings.Split(doc.Find("body").Text(), `@`)[1])
	return bibtex, nil
}
