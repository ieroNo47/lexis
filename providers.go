// providers defines the answer providers interface and implementations
package main

import (
	"math/rand"
	"time"
)

// answerProvider is an interface that defines a method to get the answer for the game
type answerProvider interface {
	getAnswer() string
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

type randomAnswerProvider struct {
	answers []string
}

func (p randomAnswerProvider) getAnswer() string {
	time.Sleep(5 * time.Second) // simulate loading time
	return p.answers[rand.Intn(len(p.answers))]
}

func newRandomAnswerProvider() randomAnswerProvider {
	return randomAnswerProvider{
		answers: []string{"apple", "grape", "peach", "mango", "berry", "lemon", "pearl"},
	}
}
