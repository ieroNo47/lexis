// providers defines the answer providers interface and implementations
package main

import (
	"math/rand"
	"slices"
)

// answerProvider is an interface that defines a method to get the answer for the game
type answerProvider interface {
	init()
	getAnswer() string
	validWord(word string) bool
}

// staticAnswerProvider is a simple implementation of answerProvider that returns a static answer
type staticAnswerProvider struct {
	answer string
}

func (p staticAnswerProvider) getAnswer() string {
	return p.answer
}

func newStaticAnswerProvider() staticAnswerProvider {
	return staticAnswerProvider{
		answer: "vogue",
	}
}

// randomAnswerProvider is an implementation of answerProvider that returns a random answer from a static list of answers
type randomAnswerProvider struct {
	answer string
	words  []string
}

func (p *randomAnswerProvider) init() {
	// time.Sleep(2 * time.Second) // simulate loading time
	p.answer = p.words[rand.Intn(len(p.words))]
}

func (p randomAnswerProvider) getAnswer() string {
	return p.answer
}

func (p randomAnswerProvider) validWord(word string) bool {
	// time.Sleep(2 * time.Second) // simulate delay
	return slices.Contains(p.words, word)
}

func newRandomAnswerProvider() *randomAnswerProvider {
	return &randomAnswerProvider{
		words: []string{"apple", "grape", "peach", "mango", "berry", "lemon", "pearl"},
	}
}
