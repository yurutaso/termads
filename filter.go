package main

type Filter struct {
	linktypes string
}

func NewFilter() *Filter {
	return &Filter{}
}

func (filter *Filter) Set(linktypes string) {
	filter.linktypes = linktypes
	return
}
