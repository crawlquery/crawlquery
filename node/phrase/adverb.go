package phrase

var AdverbialPhraseTemplates = PhraseSubCategory{
	// just a few months ago
	{"RB", "DT", "JJ", "NNS", "RB"},
	// examples of other adverbial phrases
	{"RB", "IN", "NN"},       // soon after dinner
	{"RB", "IN", "DT", "NN"}, // right after the meeting
}

func adverbialPhraseSubCategories() PhraseSubCategories {
	subCategories := PhraseSubCategories{
		"adverbial": AdverbialPhraseTemplates,
	}

	return subCategories
}
