package quizsteps

import "context"

type InMemoryQuestionsGetter struct{}

func (iqh *InMemoryQuestionsGetter) Get(context.Context) ([][]string, error) {
	return [][]string{
		{"What is the capital of France?", "Paris"},
		{"What is the capital of Ireland?", "Dublin"},
		{"What is the capital of Germany?", "Berlin"},
	}, nil
}
