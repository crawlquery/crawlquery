package phrase

var VerbTemplates = PhraseSubCategory{
	{"VBD"},
	{"VBZ"},
	{"VBG"},
	{"VBP"},
}

var VerbAdverbTemplates = PhraseSubCategory{
	{"VBD", "RB"},
	{"VBZ", "RB"},
	{"VBG", "RB"},
	{"VBP", "RB"},
}

func verbPhraseSubCategories() PhraseSubCategories {
	subCategories := PhraseSubCategories{
		"verb": VerbTemplates,
	}

	return subCategories
}
