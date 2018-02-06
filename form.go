package main

import (
	"fmt"
	"net/url"
)

// ADSform
type Form struct {
	keys   []string
	values url.Values
}

func NewForm() *Form {
	keys := []string{`db_key`, `qform`, `arxiv_sel`, `sim_query`, `ned_query`, `adsobj_query`, `aut_logic`, `obj_logic`, `nr_to_return`, `start_nr`, `jou_pick`, `ref_stems`, `data_and`, `group_and`, `start_entry_day`, `start_entry_mon`, `start_entry_year`, `end_entry_day`, `end_entry_mon`, `end_entry_year`, `min_score`, `sort`, `data_type`, `aut_syn`, `ttl_syn`, `txt_syn`, `aut_wt`, `obj_wt`, `ttl_wt`, `txt_wt`, `aut_wgt`, `obj_wgt`, `ttl_wgt`, `txt_wgt`, `ttl_sco`, `txt_sco`, `version`, `author`, `object`, `start_mon`, `start_year`, `end_mon`, `end_year`, `text`, `ttl_req`, `ttl_logic`, `title`, `txt_logic`}
	values := url.Values{}
	values.Add(`db_key`, `AST`)
	values.Add(`db_key`, `PRE`)
	values.Add(`qform`, `AST`)
	values.Add(`arxiv_sel`, `astro-ph`)
	values.Add(`arxiv_sel`, `cond-mat`)
	values.Add(`arxiv_sel`, `cs`)
	values.Add(`arxiv_sel`, `gr-qc`)
	values.Add(`arxiv_sel`, `hep-ex`)
	values.Add(`arxiv_sel`, `hep-lat`)
	values.Add(`arxiv_sel`, `hep-ph`)
	values.Add(`arxiv_sel`, `hep-th`)
	values.Add(`arxiv_sel`, `math`)
	values.Add(`arxiv_sel`, `math-ph`)
	values.Add(`arxiv_sel`, `nlin`)
	values.Add(`arxiv_sel`, `nucl-ex`)
	values.Add(`arxiv_sel`, `nucl-th`)
	values.Add(`arxiv_sel`, `physics`)
	values.Add(`arxiv_sel`, `quant-ph`)
	values.Add(`arxiv_sel`, `q-bio`)
	values.Add(`sim_query`, `YES`)
	values.Add(`ned_query`, `YES`)
	values.Add(`adsobj_query`, `YES`)
	values.Add(`nr_to_return`, `200`)
	values.Add(`start_nr`, `1`)
	values.Add(`jou_pick`, `ALL`)
	values.Add(`ref_stems`, ``)
	values.Add(`data_and`, `ALL`)
	values.Add(`group_and`, `ALL`)
	values.Add(`start_entry_day`, ``)
	values.Add(`start_entry_mon`, ``)
	values.Add(`start_entry_year`, ``)
	values.Add(`end_entry_day`, ``)
	values.Add(`end_entry_mon`, ``)
	values.Add(`end_entry_year`, ``)
	values.Add(`min_score`, ``)
	values.Add(`sort`, `SCORE`)
	values.Add(`data_type`, `SHORT`)
	values.Add(`aut_syn`, `YES`)
	values.Add(`ttl_syn`, `YES`)
	values.Add(`txt_syn`, `YES`)
	values.Add(`aut_wt`, `1.0`)
	values.Add(`obj_wt`, `1.0`)
	values.Add(`ttl_wt`, `0.3`)
	values.Add(`txt_wt`, `3.0`)
	values.Add(`aut_wgt`, `YES`)
	values.Add(`obj_wgt`, `YES`)
	values.Add(`ttl_wgt`, `YES`)
	values.Add(`txt_wgt`, `YES`)
	values.Add(`ttl_sco`, `YES`)
	values.Add(`txt_sco`, `YES`)
	values.Add(`version`, `1`)
	// Important
	values.Add(`author`, ``)
	values.Add(`aut_logic`, `OR`)
	//values.Add(`aut_req`, `YES`)
	values.Add(`object`, ``)
	values.Add(`obj_logic`, `OR`)
	values.Add(`start_mon`, ``)
	values.Add(`start_year`, ``)
	values.Add(`end_mon`, ``)
	values.Add(`end_year`, ``)
	values.Add(`title`, ``)
	values.Add(`ttl_logic`, `OR`)
	values.Add(`text`, ``)
	values.Add(`txt_logic`, `OR`)
	values.Add(`ttl_req`, `YES`)
	return &Form{keys: keys, values: values}
}

func (form *Form) Has(key string) bool {
	for _, k := range form.keys {
		if k == key {
			return true
		}
	}
	return false
}

func (form *Form) Add(key string, val string) error {
	if form.Has(key) {
		form.values.Add(key, val)
		return nil
	} else {
		return fmt.Errorf(`Fail to add value`)
	}
}

func (form *Form) AddAuthor(val string) error {
	return form.Add(`author`, val)
}

func (form *Form) AddStartDate(year string, month string) error {
	err := form.Add(`start_year`, year)
	if err != nil {
		return err
	} else {
		return form.Add(`start_mon`, month)
	}
}

func (form *Form) AddEndDate(year string, month string) error {
	err := form.Add(`end_year`, year)
	if err != nil {
		return err
	} else {
		return form.Add(`end_mon`, month)
	}
}

func (form *Form) AddTitle(val string) error {
	return form.Add(`title`, val)
}

func (form *Form) AddText(val string) error {
	return form.Add(`text`, val)
}

func (form *Form) Set(key string, val string) error {
	if form.Has(key) {
		form.values.Set(key, val)
		return nil
	} else {
		return fmt.Errorf(`Fail to set value`)
	}
}

func (form *Form) SetAuthor(val string) error {
	return form.Set(`author`, val)
}

func (form *Form) SetStartDate(year string, month string) error {
	err := form.Set(`start_year`, year)
	if err != nil {
		return err
	} else {
		return form.Set(`start_mon`, month)
	}
}

func (form *Form) SetEndDate(year string, month string) error {
	err := form.Set(`end_year`, year)
	if err != nil {
		return err
	} else {
		return form.Set(`end_mon`, month)
	}
}

func (form *Form) SetTitle(val string) error {
	return form.Set(`title`, val)
}

func (form *Form) SetText(val string) error {
	return form.Set(`text`, val)
}

func (form *Form) SetSearchLogic(key string, method string) error {
	if method == `AND` || method == `OR` {
		switch key {
		case `author`:
			form.Set(`aut_logic`, method)
			return nil
		case `title`:
			form.Set(`ttl_logic`, method)
			return nil
		case `text`:
			form.Set(`txt_logic`, method)
			return nil
		case `object`:
			form.Set(`obj_logic`, method)
			return nil
		case `all`:
			form.Set(`aut_logic`, method)
			form.Set(`ttl_logic`, method)
			form.Set(`txt_logic`, method)
			form.Set(`obj_logic`, method)
			return nil
		default:
			return fmt.Errorf(`key must be "author", "title", "text", "object" or "all"`)
		}
	} else {
		return fmt.Errorf(`method must be "AND" or "OR"`)
	}
}
