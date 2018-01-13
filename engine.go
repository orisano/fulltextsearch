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
