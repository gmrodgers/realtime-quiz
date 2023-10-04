package quizsteps

import (
	"context"
	"encoding/json"
	"fmt"
)

type QueueResultsBroadcaster struct {
	Broadcaster QueueBroadcaster
}

type quizLeaderboard struct {
	Type        string         `json:"type"`
	Leaderboard map[string]int `json:"leaderboard"`
}

func (qrb *QueueResultsBroadcaster) Broadcast(ctx context.Context, leaderboard map[string]int) error {
	body, err := json.Marshal(quizLeaderboard{Type: "leaderboard", Leaderboard: leaderboard})
	if err != nil {
		return fmt.Errorf("Failed to marshal leaderboard: '%v'", err)
	}

	err = qrb.Broadcaster.Broadcast(ctx, body)
	if err != nil {
		return fmt.Errorf("Failed to send leaderboard: '%v'", err)
	}

	return nil
}
