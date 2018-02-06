package termads

func Find(papers []*Paper, filter *Filter) [][4]string {
	result := make([][4]string, 0, len(papers))
	for _, paper := range papers {
		linktypes := paper.AvailableLinkTypesIn(filter.linktypes)
		if linktypes != "" {
			result = append(result, [4]string{paper.authors, paper.title, paper.bibcode, linktypes})
		}
	}
	return result
}
