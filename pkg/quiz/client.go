package quiz

import (
	"context"
	"fmt"
	"log"
)

type PlayerRegisterer interface {
	Register(context.Context) (bool, error)
}

type ResultsChecker interface {
	Check(context.Context) (map[string]int, error)
}

type QuestionReceiver interface {
	Receive(context.Context) (string, bool, error)
}

type AnswerFinder interface {
	Find(context.Context, string) (string, error)
}

type AnswerSender interface {
	Send(context.Context, string) error
}

type Client struct {
	PlayerRegisterer PlayerRegisterer
	QuestionReceiver QuestionReceiver
	AnswerFinder     AnswerFinder
	AnswerSender     AnswerSender
	ResultsChecker   ResultsChecker
}

func (c *Client) Join(ctx context.Context) error {
	joined, err := c.PlayerRegisterer.Register(ctx)
	if err != nil {
		return fmt.Errorf("Failed to register: '%v'", err)
	}

	if !joined {
		return nil
	}

	for {
		question, final, err := c.QuestionReceiver.Receive(ctx)
		if err != nil {
			return fmt.Errorf("Failed to receive question: '%v'", err)
		}
		log.Println(question)

		answer, err := c.AnswerFinder.Find(ctx, question)
		if err != nil {
			return fmt.Errorf("Failed to find answer: '%v'", err)
		}
		log.Println(answer)

		err = c.AnswerSender.Send(ctx, answer)
		if err != nil {
			return fmt.Errorf("Failed to send answer: '%v'", err)
		}

		if final {
			leaderboard, err := c.ResultsChecker.Check(ctx)
			if err != nil {
				return fmt.Errorf("Failed to check leaderboard: '%v'", err)
			}
			if leaderboard != nil {
				log.Println("Finished")
				log.Println(leaderboard)
				break
			}
		}
	}

	return nil
}
