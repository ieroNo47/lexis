// providers defines the answer providers interface and implementations
package main

import (
	"math/rand"
	"slices"
	"time"
)

// answerProvider is an interface that defines a method to get the answer for the game
type answerProvider interface {
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
	answers []string
}

func (p randomAnswerProvider) getAnswer() string {
	time.Sleep(1 * time.Second) // simulate loading time
	return p.answers[rand.Intn(len(p.answers))]
}

func (p randomAnswerProvider) validWord(word string) bool {
	return slices.Contains(p.answers, word)
}

func newRandomAnswerProvider() randomAnswerProvider {
	return randomAnswerProvider{
		answers: []string{"apple", "grape", "peach", "mango", "berry", "lemon", "pearl"},
	}
}
