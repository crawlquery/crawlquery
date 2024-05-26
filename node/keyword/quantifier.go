package keyword

var QuantifierNounTemplates = KeywordSubCategory{
	// few holes
	{"JJ", "NNS"},
}

func quantifierKeywordSubCategories() KeywordSubCategories {
	subCategories := KeywordSubCategories{
		"quantifier_noun": QuantifierNounTemplates,
	}

	return subCategories
}
