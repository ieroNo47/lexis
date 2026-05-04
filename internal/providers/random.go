// providers
package providers

import (
	"math/rand"
	"slices"
)

// randomAnswerProvider is an implementation of answerProvider that returns a random answer from a static list of answers
type randomAnswerProvider struct {
	answer string
	words  []string
}

func (p *randomAnswerProvider) Init() {
	// time.Sleep(2 * time.Second) // simulate loading time
	p.answer = p.words[rand.Intn(len(p.words))]
}

func (p randomAnswerProvider) GetAnswer() string {
	return p.answer
}

func (p randomAnswerProvider) ValidWord(word string) bool {
	// time.Sleep(2 * time.Second) // simulate delay
	return slices.Contains(p.words, word)
}

func NewRandomAnswerProvider() *randomAnswerProvider {
	return &randomAnswerProvider{
		words: []string{"apple", "grape", "peach", "mango", "berry", "lemon", "pearl"},
	}
}
