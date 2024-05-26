package keyword

var SimpleNounTemplates = KeywordSubCategory{
	// A tree
	{"NN"},

	// A search engine
	{"NN", "NN"},
}

var AdjectiveNounTemplates = KeywordSubCategory{
	// bright red car
	{"JJ", "JJ", "NN"},

	// lazy dog
	{"JJ", "NN"},

	//quick brown fox
	{"JJ", "NN", "NN"},

	// Best way to detect bot from user agent
	{"JJS", "NN", "TO", "VB", "NN", "IN", "NN", "NN"},
	{"JJS", "NN", "TO", "VB", "NN", "IN", "JJ", "NN"},
}

func nounKeywordSubCategories() KeywordSubCategories {
	subCategories := KeywordSubCategories{
		"simple_noun":    SimpleNounTemplates,
		"adjective_noun": AdjectiveNounTemplates,
	}

	return subCategories
}
