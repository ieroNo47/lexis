// providers defines the answer providers interface and implementations
package providers

// answerProvider is an interface that defines a method to get the answer for the game
type AnswerProvider interface {
	Init()
	GetAnswer() string
	ValidWord(word string) bool
}
