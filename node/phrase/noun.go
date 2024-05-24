package phrase

var SimpleNounTemplates = PhraseSubCategory{
	// A tree
	{"NN"},

	// A search engine
	{"NN", "NN"},
}

var AdjectiveNounTemplates = PhraseSubCategory{
	// bright red car
	{"JJ", "JJ", "NN"},

	// lazy dog
	{"JJ", "NN"},

	//quick brown fox
	{"JJ", "NN", "NN"},
}

func nounPhraseSubCategories() PhraseSubCategories {
	subCategories := PhraseSubCategories{
		"simple_noun":    SimpleNounTemplates,
		"adjective_noun": AdjectiveNounTemplates,
	}

	return subCategories
}
