package keyword

var AdverbialKeywordTemplates = KeywordSubCategory{
	// just a few months ago
	{"RB", "DT", "JJ", "NNS", "RB"},
	// examples of other adverbial keywords
	{"RB", "IN", "NN"},       // soon after dinner
	{"RB", "IN", "DT", "NN"}, // right after the meeting
}

func adverbialKeywordSubCategories() KeywordSubCategories {
	subCategories := KeywordSubCategories{
		"adverbial": AdverbialKeywordTemplates,
	}

	return subCategories
}
