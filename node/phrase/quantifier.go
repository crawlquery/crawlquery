package phrase

var QuantifierNounTemplates = PhraseSubCategory{
	// few holes
	{"JJ", "NNS"},
}

func quantifierPhraseSubCategories() PhraseSubCategories {
	subCategories := PhraseSubCategories{
		"quantifier_noun": QuantifierNounTemplates,
	}

	return subCategories
}
