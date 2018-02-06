package main

import (
	"fmt"
	"github.com/yurutaso/termads"
	"log"
)

func main() {
	// values to post
	form := termads.NewForm()
	form.SetTitle(``)
	form.SetAuthor(`tanaka`)
	form.SetStartDate(`2014`, `1`)
	form.SetEndDate(`2014`, `12`)
	form.SetSearchLogic(`all`, `AND`)

	// get links and bibcodes from doc
	papers, err := termads.GetLinks(form)
	if err != nil {
		log.Fatal(err)
	}
	for _, paper := range papers {
		abs, err := termads.GetAbstract(paper)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(abs)
	}
	//bibtex, err := GetBibTexEntry(papers[0].bibcode)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(bibtex)
	//filter := NewFilter()
	//filter.Set("FE")
	//result := Find(papers, filter)
	//fmt.Println(result)
	return
}
