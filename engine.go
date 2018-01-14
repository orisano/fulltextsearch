package fulltextsearch

import (
	"sync"
)

type Indexer interface {
	Index() string
}

type Posting struct {
	DocumentID int
	Offset     int
	Length     int
}

type Engine struct {
	tokenizer Tokenizer

	transposeIndex map[string][]Posting
	documents      []Indexer
	currentID      int
	mu             sync.RWMutex
}

func NewEngine(tokenizer Tokenizer) *Engine {
	return &Engine{
		tokenizer:      tokenizer,
		transposeIndex: map[string][]Posting{},
	}
}

func (e *Engine) AddDocument(indexer Indexer) int {
	tokens := e.tokenizer.Tokenize(indexer.Index())

	e.mu.Lock()
	defer e.mu.Unlock()
	e.documents = append(e.documents, indexer)
	id := e.currentID
	e.currentID++
	for _, token := range tokens {
		e.transposeIndex[token.Text] = append(e.transposeIndex[token.Text], Posting{
			DocumentID: id,
			Offset:     token.Offset,
			Length:     len(token.Text),
		})
	}
	return id
}

func (e *Engine) SearchOne(query string) []Posting {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.transposeIndex[query]
}

func and(a, b []Posting) []Posting {
	var p []Posting

	both := map[int]struct{}{}
	for _, x := range a {
		for len(b) > 0 {
			docID := b[0].DocumentID
			if docID > x.DocumentID {
				break
			}
			if docID == x.DocumentID {
				both[docID] = struct{}{}
				if b[0].Offset > x.Offset {
					break
				}
			}
			if _, ok := both[docID]; ok {
				p = append(p, b[0])
			}
			b = b[1:]
		}
		if _, ok := both[x.DocumentID]; ok {
			p = append(p, x)
		}
	}
	p = append(p, b...)
	return p
}

func (e *Engine) SearchAnd(queries []string) []Posting {
	if len(queries) == 0 {
		return nil
	}

	e.mu.RLock()
	defer e.mu.RUnlock()
	r := e.transposeIndex[queries[0]]
	for _, query := range queries[1:] {
		r = and(r, e.transposeIndex[query])
	}
	return r
}

func merge(a, b []Posting) []Posting {
	var p []Posting

	for _, x := range a {
		for len(b) > 0 {
			docID := b[0].DocumentID
			if docID > x.DocumentID {
				break
			}
			if docID == x.DocumentID {
				if b[0].Offset > x.Offset {
					break
				}
			}
			p = append(p, b[0])
			b = b[1:]
		}
		p = append(p, x)
	}
	p = append(p, b...)
	return p
}

func (e *Engine) SearchOr(queries []string) []Posting {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var p []Posting
	for _, query := range queries {
		p = merge(p, e.transposeIndex[query])
	}
	return p
}
