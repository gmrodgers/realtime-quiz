package quizsteps

import (
	"context"
	"encoding/json"
	"fmt"
)

type QueueQuestionReceiver struct {
	Listener QueueListener
}

func (qqr *QueueQuestionReceiver) Receive(ctx context.Context) (string, bool, error) {
	body, err := qqr.Listener.Listen(ctx, func(b [][]byte) bool {
		var question quizQuestion
		err := json.Unmarshal(b[len(b)-1], &question)
		if err != nil {
			return false
		}
		return question.Type == "question"
	})
	if err != nil {
		return "", false, fmt.Errorf("Failed to receive question: '%v'", err)
	}

	var question quizQuestion
	err = json.Unmarshal(body[len(body)-1], &question)
	if err != nil {
		return "", false, fmt.Errorf("Failed to unmarshal question: '%v'", err)
	}

	return question.Question, question.Final, nil
}
