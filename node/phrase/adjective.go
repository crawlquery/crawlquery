package phrase

var AdjectivePhraseTemplates = PhraseSubCategory{
	{"RB", "JJ"},
	{"JJ"},
	{"JJS"},
}

func adjectivePhraseSubCategories() PhraseSubCategories {
	subCategories := PhraseSubCategories{
		"adjective": AdjectivePhraseTemplates,
	}

	return subCategories
}
