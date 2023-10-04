package quizsteps

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type QueueMessager interface {
	Send(ctx context.Context, msg []byte) error
	Receive(ctx context.Context, until func([][]byte) bool) ([][]byte, error)
}

type QueuePlayerRegistrar struct {
	Name     string
	Messager QueueMessager
}

type quizMessage struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func (qpr *QueuePlayerRegistrar) StartRegistration(ctx context.Context, n int, t time.Duration) ([]string, error) {
	ctxTimeout, cancelFunc := context.WithTimeout(ctx, t)
	defer cancelFunc()
	resps, err := qpr.Messager.Receive(ctxTimeout, func(b [][]byte) bool {
		log.Println("...")
		return len(b) == n
	})
	if err != nil {
		return []string{}, fmt.Errorf("Failed to receive: '%v'", err)
	}

	log.Println("players are:")
	var registrations []string
	for _, resp := range resps {
		var qm quizMessage
		err = json.Unmarshal(resp, &qm)
		if err != nil {
			continue
		}
		log.Println(qm.Payload)
		registrations = append(registrations, qm.Payload)
	}

	// TODO: more just that we got a message from them
	// Should send an ack back to tell them they've got the place!
	log.Println("all players registered")

	return registrations, nil
}

func (rpr *QueuePlayerRegistrar) Register(ctx context.Context) (bool, error) {
	body, err := json.Marshal(quizMessage{Type: "join", Payload: rpr.Name})
	if err != nil {
		return false, fmt.Errorf("Failed to marshal message: '%v'", err)
	}

	err = rpr.Messager.Send(ctx, body)
	if err != nil {
		return false, fmt.Errorf("Faild to send registration message: '%v'", err)
	}

	return true, nil
}
