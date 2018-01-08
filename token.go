package fulltextsearch

type Tokenizer interface {
	Tokenize(string) []string
}

type NgramTokenizer struct {
	N int
}

func (nt *NgramTokenizer) Tokenize(s string) []string {
	r := []rune(s)
	t := len(r) - nt.N + 1
	if t <= 0 {
		return nil
	}
	tokens := make([]string, 0, t)
	for i := 0; i < t; i++ {
		tokens = append(tokens, string(r[i:i+nt.N]))
	}
	return tokens
}
