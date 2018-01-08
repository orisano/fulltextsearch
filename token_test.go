package fulltextsearch

import "testing"

func equals(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestNgramTokenizer_Tokenize(t *testing.T) {
	tests := []struct {
		n        int
		text     string
		expected []string
	}{
		{
			n:        3,
			text:     "foobar",
			expected: []string{"foo", "oob", "oba", "bar"},
		},
		{
			n:        2,
			text:     "検索エンジン",
			expected: []string{"検索", "索エ", "エン", "ンジ", "ジン"},
		},
		{
			n:        2,
			text:     "𩸽の定食",
			expected: []string{"𩸽の", "の定", "定食"},
		},
		{
			n:        10,
			text:     "less",
			expected: []string{},
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
