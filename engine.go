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
}

type Engine struct {
	tokenizer Tokenizer

	transposeIndex map[string][]Posting
	documents      []Indexer
	currentID      int
	mu             sync.RWMutex
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
		})
	}
	return id
}
