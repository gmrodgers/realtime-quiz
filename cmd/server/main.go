package main

import (
	"context"
	"flag"
	"log"

	"github.com/gmrodgers/realtime-quiz/pkg/quiz"
	"github.com/gmrodgers/realtime-quiz/pkg/quizsteps"
	"github.com/gmrodgers/realtime-quiz/pkg/repositories/rabbitmq"
)

const (
	host = "amqp://guest:guest@localhost:5672/"
)

func main() {
	var numberOfPlayers int
	flag.IntVar(&numberOfPlayers, "n", 0, "number of players")
	flag.Parse()
	quiz := quiz.Server{
		PlayerRegistrar: &quizsteps.QueuePlayerRegistrar{
			Messager: &rabbitmq.DirectQueue{Host: host, Queue: "registration"},
		},
		QuestionsGetter: &quizsteps.InMemoryQuestionsGetter{},
		QuestionBroadcaster: &quizsteps.QueueQuestionBroadcaster{
			Broadcaster: &rabbitmq.FanoutQueue{Host: host, Exchange: "questions"},
			Messager:    &rabbitmq.DirectQueue{Host: host, Queue: "answers"},
		},
		ResultsBroadcaster: &quizsteps.QueueResultsBroadcaster{
			Broadcaster: &rabbitmq.FanoutQueue{Host: host, Exchange: "results"},
		},
	}

	if err := quiz.Host(context.Background(), numberOfPlayers); err != nil {
		log.Fatalln(err)
	}
}
