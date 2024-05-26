package keyword

import (
	"reflect"
	"testing"

	"github.com/jdkato/prose/v2"
)

func TestParseQuantifierNounKeywords(t *testing.T) {
	QuantifierNounTemplates := KeywordSubCategory{
		// few holes
		{"JJ", "NNS"},
	}

	subCategories := KeywordSubCategories{
		"quantifier_noun": QuantifierNounTemplates,
	}

	cases := []struct {
		name     string
		sentence string
		want     [][]string
	}{
		{
			name:     "Quantifier noun keyword",
			sentence: "There are a few holes in the bucket.",
			want:     [][]string{{"few", "holes"}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := prose.NewDocument(tc.sentence)
			if err != nil {
				t.Fatalf("Failed to parse document: %v", err)
			}

			got := parseSubCategories(doc.Tokens(), subCategories)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Expected %v, got %v", tc.want, got)
			}
		})
	}
}
