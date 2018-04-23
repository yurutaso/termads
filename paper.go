package termads

import (
	"fmt"
	"strings"
)

/*=======================================================
/*                    Paper
/*=======================================================*/
const (
	LINKTYPE_ABSTRACT           string = `A`
	LINKTYPE_CITATIONS          string = `C`
	LINKTYPE_ONLINE_DATA        string = `D`
	LINKTYPE_ELEC_ARTICLE       string = `E`
	LINKTYPE_FULL_ARTICLE       string = `F`
	LINKTYPE_GIF                string = `G`
	LINKTYPE_HEP                string = `H`
	LINKTYPE_ADDITIONAL_INFO    string = `I`
	LINKTYPE_LIBRARY_ENTRIES    string = `L`
	LINKTYPE_MULTIMEDIA         string = `M`
	LINKTYPE_NED                string = `N`
	LINKTYPE_ASSOCIATED_ARTICLE string = `O`
	LINKTYPE_PLANETARY_DATA     string = `P`
	LINKTYPE_REFERENCES         string = `R`
	LINKTYPE_SIMBAD             string = `S`
	LINKTYPE_TABLE_OF_CONTENTS  string = `T`
	LINKTYPE_ALSO_READ_ARTICLE  string = `U`
	LINKTYPE_ARXIV              string = `X`
	LINKTYPE_ABSTRACT_CUSTOM    string = `Z`
	VALID_LINKS                 string = `ACDEFGHILMNOPRSTUXZ`
)

type Paper interface {
	String() string
	GetBibTex() (string, error)
	GetBibcode() string
	SetBibcode(string)
	GetURLOfType(string) string
	SetURL(string, string) error
	GetTitle() string
	SetTitle(string)
	GetAuthors() string
	SetAuthors(string)
	GetAbstract() string
	SetAbstract(string)
	SetAbstractFromADS() error
	LinkTypes() string
	LinkTypesIn(string) string
	HasLink(string) bool
	HasLinkAll(string) bool
}

type paper struct {
	bibcode  string
	title    string
	authors  string
	abstract string
	links    map[string]string
}

func NewPaper() Paper {
	return &paper{
		bibcode: "",
		links:   make(map[string]string)}
}

func (p *paper) String() string {
	s := ""
	for key, val := range p.links {
		if val != "" {
			s += fmt.Sprintf("%s : %s", key, val)
		}
	}
	return s
}

func (p *paper) GetURLOfType(linktype string) string {
	if p.HasLink(linktype) {
		return p.links[linktype]
	}
	return ""
}

func (p *paper) GetBibcode() string {
	return p.bibcode
}

func (p *paper) SetBibcode(bibcode string) {
	p.bibcode = bibcode
}
func (p *paper) GetBibTex() (string, error) {
	return GetBibTex(p.bibcode)
}

func (p *paper) GetAuthors() string {
	return p.authors
}

func (p *paper) SetAuthors(authors string) {
	p.authors = authors
}

func (p *paper) GetAbstract() string {
	return p.abstract
}

func (p *paper) SetAbstract(abstract string) {
	p.abstract = abstract
}

func (p *paper) SetAbstractFromADS() error {
	if p.HasLink(LINKTYPE_ABSTRACT) {
		abs, err := GetAbstract(p.links[LINKTYPE_ABSTRACT])
		if err != nil {
			return err
		}
		p.abstract = abs
		return nil
	}
	return fmt.Errorf(`this paper does not have linktype %s`, LINKTYPE_ABSTRACT)
}

func (p *paper) GetTitle() string {
	return p.title
}

func (p *paper) SetTitle(title string) {
	p.title = title
}

func (p *paper) SetURL(url, linktype string) error {
	if len(linktype) == 1 && strings.Contains(VALID_LINKS, linktype) {
		p.links[strings.ToUpper(linktype)] = url
		return nil
	}
	return fmt.Errorf(`linktype must be a single character in %v`, VALID_LINKS)
}

func (p *paper) LinkTypes() string {
	linktypes := ""
	for _, linktype := range VALID_LINKS {
		s := string(linktype)
		if p.HasLink(s) {
			linktypes += s
		}
	}
	return linktypes
}

func (p *paper) LinkTypesIn(linktypes string) string {
	_linktypes := ""
	for _, linktype := range VALID_LINKS {
		s := string(linktype)
		if p.HasLink(s) && strings.Contains(linktypes, s) {
			_linktypes += s
		}
	}
	return _linktypes
}

func (p *paper) HasLink(linktypes string) bool {
	for _, linktype := range linktypes {
		if p.HasLink(string(linktype)) {
			return true
		}
	}
	return false
}

func (p *paper) HasLinkAll(linktype string) bool {
	linktype = strings.ToUpper(linktype)
	if len(linktype) == 1 && strings.Contains(VALID_LINKS, linktype) {
		if p.links[linktype] == "" {
			return false
		} else {
			return true
		}
	}
	fmt.Println(`linktype must be a single character in %v`, VALID_LINKS)
	return false
}
