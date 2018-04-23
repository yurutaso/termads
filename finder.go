package termads

func Find(papers []Paper, filter *Filter) [][4]string {
	result := make([][4]string, 0, len(papers))
	for _, paper := range papers {
		linktypes := paper.LinkTypesIn(filter.linktypes)
		if linktypes != "" {
			result = append(result, [4]string{paper.GetAuthors(), paper.GetTitle(), paper.GetBibcode(), linktypes})
		}
	}
	return result
}
