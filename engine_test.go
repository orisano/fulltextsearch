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
				{0, 2, 3},
				{2, 0, 3},
				{5, 1, 3},
			},
		},
		{
			engine: buildEngine(&NgramTokenizer{3}, []string{
				"example", "日本語", "amplify", "foo", "bar", "campfire",
			}),
			query: "日本語",
			expected: []Posting{
				{1, 0, len("日本語")},
			},
		},
	}

	for _, test := range tests {
		if got := test.engine.SearchOne(test.query); !reflect.DeepEqual(got, test.expected) {
			t.Errorf("unexpected posting list. expected: %#v, but got: %#v", test.expected, got)
		}
	}
}

func TestEngine_SearchAnd(t *testing.T) {
	tests := []struct {
		engine   *Engine
		queries  []string
		expected []Posting
	}{
		{
			engine: buildEngine(&NgramTokenizer{3}, []string{
				"foo", "a foo bar example amplify", "amplify fool bare", "foo camp", "bar bar", "campfire",
			}),
			queries: []string{"foo", "amp", "bar"},
			expected: []Posting{
				{1, 2, 3},
				{1, 6, 3},
				{1, 12, 3},
				{1, 18, 3},
				{2, 0, 3},
				{2, 8, 3},
				{2, 13, 3},
			},
		},
	}

	for _, test := range tests {
		if got := test.engine.SearchAnd(test.queries); !reflect.DeepEqual(got, test.expected) {
			t.Errorf("unexpected posting list. expected: %#v, but got: %#v", test.expected, got)
		}
	}
}

func TestEngine_SearchOr(t *testing.T) {
	tests := []struct {
		engine   *Engine
		queries  []string
		expected []Posting
	}{
		{
			engine: buildEngine(&NgramTokenizer{3}, []string{
				"foo", "a foo bar example amplify", "amplify fool bare", "foo camp", "bar bar", "campfire",
			}),
			queries: []string{"foo", "amp", "bar"},
			expected: []Posting{
				{0, 0, 3},
				{1, 2, 3},
				{1, 6, 3},
				{1, 12, 3},
				{1, 18, 3},
				{2, 0, 3},
				{2, 8, 3},
				{2, 13, 3},
				{3, 0, 3},
				{3, 5, 3},
				{4, 0, 3},
				{4, 4, 3},
				{5, 1, 3},
			},
		},
	}

	for _, test := range tests {
		if got := test.engine.SearchOr(test.queries); !reflect.DeepEqual(got, test.expected) {
			t.Errorf("unexpected posting list. expected: %#v, but got: %#v", test.expected, got)
		}
	}
}
