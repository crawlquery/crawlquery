package keyword

var AdjectiveKeywordTemplates = KeywordSubCategory{
	{"RB", "JJ"},
	{"JJ"},
	{"JJS"},
}

func adjectiveKeywordSubCategories() KeywordSubCategories {
	subCategories := KeywordSubCategories{
		"adjective": AdjectiveKeywordTemplates,
	}

	return subCategories
}
