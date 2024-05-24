package phrase

var PrepositionalPhraseTemplates = PhraseSubCategory{
	// price of eggs
	{"NN", "IN", "NNS"},

	// underneath the egg tray
	{"IN", "DT", "NN", "NN"},

	// price of the eggs
	{"NN", "IN", "DT", "NNS"},

	// at an all-time high
	{"IN", "DT", "JJ", "NN"},
	{"IN", "DT", "JJ", "JJ"},

	// other examples of prepositional phrases
	{"IN", "DT", "NN"}, // at the park
	{"IN", "DT", "NN"}, // in the park
	{"IN", "NNS"},      // with friends
	{"IN", "NNP"},      // in London
}

func prepositionalPhraseSubCategories() PhraseSubCategories {
	subCategories := PhraseSubCategories{
		"prepositional": PrepositionalPhraseTemplates,
	}

	return subCategories
}
