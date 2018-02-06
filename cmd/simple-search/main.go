package main

import (
	"fmt"
	"github.com/yurutaso/termads"
	"log"
)

func main() {
	// values to post
	form := termads.NewForm()
	form.SetTitle(`title`)
	form.SetAuthor(`author`)
	form.SetStartDate(`2017`, `1`)
	form.SetEndDate(`2017`, `12`)
	form.SetSearchLogic(`all`, `AND`)

	// get links and bibcodes from doc
	papers, err := termads.GetLinks(form)
	if err != nil {
		log.Fatal(err)
	}
	abs, err := termads.GetAbstract(papers[0])
	fmt.Println(abs)
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
