package main

import (
	"flag"
	"fmt"
	"github.com/yurutaso/termads"
	"log"
	"strconv"
)

const (
	MAXIMUM_RESULT int = 5
)

var (
	t   = flag.String("t", "", "title of the paper")
	n   = flag.Int("n", MAXIMUM_RESULT, "maximum number")
	v   = flag.Bool("v", false, "verbose (show abstract of the paper)")
	a   = flag.String("a", "", "author of the paper")
	abs = flag.String("abs", "", "abstract of the paper")
	y   = flag.Int("y", 0, "year")
	m   = flag.Int("m", -1, "month")
	y1  = flag.Int("y1", 0, "year begin (ignored if -y is set)")
	m1  = flag.Int("m1", 0, "month begin (ignored if -m is set)")
	y2  = flag.Int("y2", 3000, "year end (ignored if -y is set)")
	m2  = flag.Int("m2", 12, "month end (ignored if -m is set)")
)

func main() {
	flag.Parse()
	if *y != -1 {
		*y1 = *y
		*y2 = *y
	}
	if *m != -1 {
		*m1 = *m
		*m2 = *m
	}

	// values to post
	form := termads.NewForm()
	form.SetTitle(*t)
	form.SetAuthor(*a)
	form.SetText(*abs)
	form.SetStartDate(strconv.Itoa(*y1), strconv.Itoa(*m1))
	form.SetEndDate(strconv.Itoa(*y2), strconv.Itoa(*m2))
	form.SetSearchLogic(`all`, `AND`)
	form.Set(`aut_req`, `YES`)
	form.Set(`txt_req`, `YES`)

	// get links and bibcodes from doc
	papers, err := termads.GetPapers(form)
	if err != nil {
		log.Fatal(err)
	}
	count := 1
	for _, paper := range papers {
		if count > *n {
			break
		}
		count++
		// abstract
		if *v {
			err := paper.SetAbstractFromADS()
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(paper.GetAbstract())
		}
		// bibtex
		bibtex, err := paper.GetBibTex()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(bibtex)
		fmt.Println("--------------------------------------------")
	}
	return
}
