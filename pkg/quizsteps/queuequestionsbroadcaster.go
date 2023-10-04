package quizsteps

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type QueueBroadcaster interface {
	Broadcast(context.Context, []byte) error
}

type QueueQuestionBroadcaster struct {
	Broadcaster QueueBroadcaster
	Messager    QueueMessager
}

type quizQuestion struct {
	Type     string `json:"type"`
	Question string `json:"question"`
	Final    bool   `json:"final"`
}

type quizAnswer struct {
	Type   string `json:"type"`
	Player string `json:"player"`
	Answer string `json:"answer"`
}

func (qb *QueueQuestionBroadcaster) Broadcast(ctx context.Context, question string, n int, t time.Duration, final bool) (map[string]string, error) {
	body, err := json.Marshal(quizQuestion{Type: "question", Final: final, Question: question})
	if err != nil {
		return map[string]string{}, fmt.Errorf("Failed to marshal message: '%v'", err)
	}

	err = qb.Broadcaster.Broadcast(ctx, body)
	if err != nil {
		return map[string]string{}, fmt.Errorf("Failed to send message: '%v'", err)
	}

	ctxTimeout, cancelFunc := context.WithTimeout(ctx, t)
	defer cancelFunc()
	resps, err := qb.Messager.Receive(ctxTimeout, func(b [][]byte) bool {
		log.Println("...")
		return len(b) == n
	})
	if err != nil {
		return map[string]string{}, fmt.Errorf("Failed to receive: '%v'", err)
	}

	answers := map[string]string{}
	for _, resp := range resps {
		var qa quizAnswer
		err = json.Unmarshal(resp, &qa)
		if err != nil {
			continue
		}
		if qa.Type == "answer" {
			answers[qa.Player] = qa.Answer
		}
	}
	<-ctxTimeout.Done()

	return answers, nil
}
