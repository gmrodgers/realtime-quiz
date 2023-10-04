package quiz

import (
	"context"
	"fmt"
	"testing"
)

var (
	fake1 = &fakeClienter{questionFinality: []bool{true}}
	fake2 = &fakeClienter{questionFinality: []bool{false, true}}
)

func TestClient_Join(t *testing.T) {
	type fields struct {
		PlayerRegisterer PlayerRegisterer
		QuestionReceiver QuestionReceiver
		AnswerFinder     AnswerFinder
		AnswerSender     AnswerSender
		ResultsChecker   ResultsChecker
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name                          string
		fields                        fields
		args                          args
		wantErr                       bool
		wantResultsCheckAfterQuestion int
	}{
		{
			name: "NoErrorsThrownFinalQuestion_SuccessChecksResults",
			fields: fields{
				PlayerRegisterer: fake1,
				QuestionReceiver: fake1,
				AnswerFinder:     fake1,
				AnswerSender:     fake1,
				ResultsChecker:   fake1,
			},
			args:                          args{ctx: context.TODO()},
			wantErr:                       false,
			wantResultsCheckAfterQuestion: 1,
		},
		{
			name: "NoErrorsThrownNonFinalQuestion_SuccessDoesntCheckResultsOnFirst",
			fields: fields{
				PlayerRegisterer: fake2,
				QuestionReceiver: fake2,
				AnswerFinder:     fake2,
				AnswerSender:     fake2,
				ResultsChecker:   fake2,
			},
			args:                          args{ctx: context.TODO()},
			wantErr:                       false,
			wantResultsCheckAfterQuestion: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				PlayerRegisterer: tt.fields.PlayerRegisterer,
				QuestionReceiver: tt.fields.QuestionReceiver,
				AnswerFinder:     tt.fields.AnswerFinder,
				AnswerSender:     tt.fields.AnswerSender,
				ResultsChecker:   tt.fields.ResultsChecker,
			}
			err := c.Join(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Join() error = %v, wantErr %v", err, tt.wantErr)
			}
			f, _ := (c.QuestionReceiver).(*fakeClienter)
			if f.question != tt.wantResultsCheckAfterQuestion {
				t.Errorf("Client.Join() checksResults = %v, wantResultsCheck %v", err, tt.wantErr)
			}
		})
	}
}

/* Test Fakes */

var _ PlayerRegisterer = &fakeClienter{}
var _ QuestionReceiver = &fakeClienter{}
var _ AnswerFinder = &fakeClienter{}
var _ AnswerSender = &fakeClienter{}
var _ ResultsChecker = &fakeClienter{}

type fakeClienter struct {
	questionFinality []bool
	question         int
}

// Register implements PlayerRegisterer.
func (*fakeClienter) Register(context.Context) (bool, error) {
	return true, nil
}

// Receive implements QuestionReceiver.
func (f *fakeClienter) Receive(context.Context) (string, bool, error) {
	fmt.Println(f.questionFinality, f.question)
	final := f.questionFinality[f.question]
	return "Question?", final, nil
}

// Find implements AnswerFinder.
func (*fakeClienter) Find(context.Context, string) (string, error) {
	return "Answer.", nil
}

// Send implements AnswerSender.
func (f *fakeClienter) Send(context.Context, string) error {
	f.question += 1
	return nil
}

// Check implements ResultsChecker.
func (f *fakeClienter) Check(context.Context) (map[string]int, error) {
	return map[string]int{"me": 1}, nil
}
