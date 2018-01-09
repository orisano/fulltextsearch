package fulltextsearch

import "unicode/utf8"

type Token struct {
	Text   string
	Offset int
}

type Tokenizer interface {
	Tokenize(string) []*Token
}

type NgramTokenizer struct {
	N int
}

func (nt *NgramTokenizer) Tokenize(s string) []*Token {
	r := []rune(s)
	t := len(r) - nt.N + 1
	if t <= 0 {
		return nil
	}
	tokens := make([]*Token, 0, t)
	offset := 0
	for i := 0; i < t; i++ {
		tokens = append(tokens, &Token{
			Text:   string(r[i : i+nt.N]),
			Offset: offset,
		})
		offset += utf8.RuneLen(r[i])
	}
	return tokens
}
