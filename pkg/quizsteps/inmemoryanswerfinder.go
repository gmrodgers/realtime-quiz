package quizsteps

import "context"

type InMemoryAnswerFinder struct{}

func (maf *InMemoryAnswerFinder) Find(ctx context.Context, question string) (string, error) {
	answers := map[string]string{
		"What is the capital of France?":  "Paris",
		"What is the capital of Ireland?": "Dublin",
		"What is the capital of Germany?": "Berlin",
	}

	return answers[question], nil
}
