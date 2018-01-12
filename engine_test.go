package fulltextsearch

import (
	"reflect"
	"testing"
)

type rawIndex string

func (r rawIndex) Index() string {
	return string(r)
}

func buildEngine(tokenizer Tokenizer, docs []string) *Engine {
	engine := NewEngine(tokenizer)
	for _, doc := range docs {
		engine.AddDocument(rawIndex(doc))
	}
	return engine
}

func TestEngine_SearchOne(t *testing.T) {
	tests := []struct {
		engine   *Engine
		query    string
		expected []Posting
	}{
		{
			engine: buildEngine(&NgramTokenizer{3}, []string{
				"example", "日本語", "amplify", "foo", "bar", "campfire",
			}),
			query: "amp",
			expected: []Posting{
				{0, 2},
				{2, 0},
				{5, 1},
			},
		},
	}

	for _, test := range tests {
		if got := test.engine.SearchOne(test.query); !reflect.DeepEqual(got, test.expected) {
			t.Errorf("unexpected posting list. expected: %#v, but got: %#v", test.expected, got)
		}
	}
}
