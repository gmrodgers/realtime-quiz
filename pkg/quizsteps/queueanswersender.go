package quizsteps

import (
	"context"
	"encoding/json"
	"fmt"
)

type QueueAnswerSender struct {
	Name     string
	Messager QueueMessager
}

func (qas *QueueAnswerSender) Send(ctx context.Context, answer string) error {
	body, err := json.Marshal(quizAnswer{Type: "answer", Player: qas.Name, Answer: answer})
	if err != nil {
		return fmt.Errorf("Failed to marshal answer: '%v'", err)
	}

	err = qas.Messager.Send(ctx, body)
	if err != nil {
		return fmt.Errorf("Failed to send answer: '%v'", err)
	}

	return nil
}
