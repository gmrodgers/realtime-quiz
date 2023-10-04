package quizsteps

import (
	"context"
	"encoding/json"
	"fmt"
)

type QueueListener interface {
	Listen(context.Context, func([][]byte) bool) ([][]byte, error)
}

type QueueResultsChecker struct {
	Listener QueueListener
}

func (qrc *QueueResultsChecker) Check(ctx context.Context) (map[string]int, error) {
	body, err := qrc.Listener.Listen(ctx, func(b [][]byte) bool {
		return len(b) == 1
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to receive results: '%v'", err)
	}
	if len(body) == 0 {
		return nil, nil
	}

	var leaderboard quizLeaderboard
	err = json.Unmarshal(body[0], &leaderboard)
	if err != nil {
		return nil, nil
		// return nil, fmt.Errorf("Failed to unmarshal results: '%v'", err)
	}

	return leaderboard.Leaderboard, nil
}
