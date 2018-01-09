package fulltextsearch

import "testing"

func equals(a, b []*Token) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Text != b[i].Text {
			return false
		}
		if a[i].Offset != b[i].Offset {
			return false
		}
	}
	return true
}

func TestNgramTokenizer_Tokenize(t *testing.T) {
	tests := []struct {
		n        int
		text     string
		expected []*Token
	}{
		{
			n:    3,
			text: "foobar",
			expected: []*Token{
				{"foo", 0},
				{"oob", 1},
				{"oba", 2},
				{"bar", 3},
			},
		},
		{
			n:    2,
			text: "検索エンジン",
			expected: []*Token{
				{"検索", 0},
				{"索エ", len("検")},
				{"エン", len("検索")},
				{"ンジ", len("検索エ")},
				{"ジン", len("検索エン")},
			},
		},
		{
			n:    2,
			text: "𩸽の定食",
			expected: []*Token{
				{"𩸽の", 0},
				{"の定", len("𩸽")},
				{"定食", len("𩸽の")},
			},
		},
		{
			n:        10,
			text:     "less",
			expected: []*Token{},
		},
	}
	for _, test := range tests {
		nt := NgramTokenizer{test.n}
		got := nt.Tokenize(test.text)
		if !equals(got, test.expected) {
			t.Errorf("unexpected tokens. expected: %v, but got: %v", test.expected, got)
		}
	}
}
