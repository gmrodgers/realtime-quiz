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
	var clientName string
	flag.StringVar(&clientName, "name", "", "name of player")
	flag.Parse()

	quiz := quiz.Client{
		PlayerRegisterer: &quizsteps.QueuePlayerRegistrar{
			Name:     clientName,
			Messager: &rabbitmq.DirectQueue{Host: host, Queue: "registration"},
		},
		QuestionReceiver: &quizsteps.QueueQuestionReceiver{
			Listener: &rabbitmq.FanoutQueue{Host: host, Exchange: "questions"},
		},
		AnswerFinder: &quizsteps.InMemoryAnswerFinder{},
		AnswerSender: &quizsteps.QueueAnswerSender{
			Name:     clientName,
			Messager: &rabbitmq.DirectQueue{Host: host, Queue: "answers"},
		},
		ResultsChecker: &quizsteps.QueueResultsChecker{
			Listener: &rabbitmq.FanoutQueue{Host: host, Exchange: "results"},
		},
	}

	if err := quiz.Join(context.Background()); err != nil {
		log.Fatalln(err)
	}
}
