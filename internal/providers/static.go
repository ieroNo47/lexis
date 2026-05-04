package providers

// staticAnswerProvider is a simple implementation of answerProvider that returns a static answer
type staticAnswerProvider struct {
	answer string
}

func (p staticAnswerProvider) GetAnswer() string {
	return p.answer
}

func (p staticAnswerProvider) ValidWord(word string) bool {
	return word == p.answer
}

func (p staticAnswerProvider) Init() {
	// do nothing
}

func NewStaticAnswerProvider() staticAnswerProvider {
	return staticAnswerProvider{
		answer: "vogue",
	}
}
