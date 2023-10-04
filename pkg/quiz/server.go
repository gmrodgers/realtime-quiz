package quiz

import (
	"context"
	"fmt"
	"log"
	"time"
)

type PlayerRegistrar interface {
	StartRegistration(context.Context, int, time.Duration) ([]string, error)
}

type QuestionsGetter interface {
	Get(context.Context) ([][]string, error)
}

type QuestionBroadcaster interface {
	Broadcast(context.Context, string, int, time.Duration, bool) (map[string]string, error)
}

type ResultsBroadcaster interface {
	Broadcast(context.Context, map[string]int) error
}

type Server struct {
	PlayerRegistrar     PlayerRegistrar
	QuestionsGetter     QuestionsGetter
	QuestionBroadcaster QuestionBroadcaster
	ResultsBroadcaster  ResultsBroadcaster
}

func (s *Server) Host(ctx context.Context, numberOfPlayers int) error {
	registrationTime := 60 * time.Second
	contestants, err := s.PlayerRegistrar.StartRegistration(ctx, numberOfPlayers, registrationTime)
	if err != nil {
		return fmt.Errorf("Failed to register: '%v'", err)
	}

	leaderboard := map[string]int{}
	for _, contestant := range contestants {
		leaderboard[contestant] = 0
	}

	qas, err := s.QuestionsGetter.Get(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get questions: '%v'", err)
	}

	// TODO: do I need this
	time.Sleep(time.Second)

	for i, qa := range qas {
		question, expectedAnswer := qa[0], qa[1]
		log.Println(question)
		questionTime := time.Second * 10
		final := i == len(qas)-1
		answers, err := s.QuestionBroadcaster.Broadcast(ctx, question, numberOfPlayers, questionTime, final)
		if err != nil {
			return fmt.Errorf("Failed to broadcast question: '%v'", err)
		}
		for contestant, answer := range answers {
			log.Println(contestant, answer)
			if answer == expectedAnswer {
				leaderboard[contestant] += 1
			}
		}
	}

	err = s.ResultsBroadcaster.Broadcast(ctx, leaderboard)
	if err != nil {
		return fmt.Errorf("Failed to broadcast results: '%v'", err)
	}
	log.Printf("Results: %v\n", leaderboard)

	return nil
}
