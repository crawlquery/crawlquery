package keyword

var VerbTemplates = KeywordSubCategory{
	{"VBD"},
	{"VBZ"},
	{"VBG"},
	{"VBP"},
	{"VB"},
}

var VerbAdverbTemplates = KeywordSubCategory{
	{"VBD", "RB"},
	{"VBZ", "RB"},
	{"VBG", "RB"},
	{"VBP", "RB"},
}

func verbKeywordSubCategories() KeywordSubCategories {
	subCategories := KeywordSubCategories{
		"verb": VerbTemplates,
	}

	return subCategories
}
