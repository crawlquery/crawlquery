package phrase

var AdjectivePhraseTemplates = PhraseSubCategory{
	{"RB", "JJ"},
}

func adjectivePhraseSubCategories() PhraseSubCategories {
	subCategories := PhraseSubCategories{
		"adjective": AdjectivePhraseTemplates,
	}

	return subCategories
}
