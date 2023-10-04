package quiz

import (
	"context"
	"testing"
	"time"
)

var (
	zeroAnswerFake  = &fakeServerer{playerAnswers: map[string]string{}}
	zeroAnswerFake2 = &fakeServerer2{}

	oneAnswerFake  = &fakeServerer{playerAnswers: map[string]string{"me": "Answer."}}
	oneAnswerFake2 = &fakeServerer2{}
)

func TestServer_Host(t *testing.T) {
	type fields struct {
		PlayerRegistrar     PlayerRegistrar
		QuestionsGetter     QuestionsGetter
		QuestionBroadcaster QuestionBroadcaster
		ResultsBroadcaster  ResultsBroadcaster
	}
	type args struct {
		ctx             context.Context
		numberOfPlayers int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "NoErrorsThrown_Success",
			fields: fields{
				PlayerRegistrar:     &fakeServerer{},
				QuestionsGetter:     &fakeServerer{},
				QuestionBroadcaster: &fakeServerer{},
				ResultsBroadcaster:  &fakeServerer2{},
			},
			args: args{
				ctx:             context.TODO(),
				numberOfPlayers: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				PlayerRegistrar:     tt.fields.PlayerRegistrar,
				QuestionsGetter:     tt.fields.QuestionsGetter,
				QuestionBroadcaster: tt.fields.QuestionBroadcaster,
				ResultsBroadcaster:  tt.fields.ResultsBroadcaster,
			}
			err := s.Host(tt.args.ctx, tt.args.numberOfPlayers)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.Host() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

/* Test Fakes */

var _ PlayerRegistrar = &fakeServerer{}
var _ QuestionsGetter = &fakeServerer{}
var _ QuestionBroadcaster = &fakeServerer{}
var _ ResultsBroadcaster = &fakeServerer2{}

type fakeServerer struct {
	playerAnswers map[string]string
}

// Needed as a fakeServerer cant be a QuestionBroadcaster and ResultsBroadcaster at the same time (method name)
type fakeServerer2 struct {
}

// StartRegistration implements PlayerRegistrar.
func (f *fakeServerer) StartRegistration(context.Context, int, time.Duration) ([]string, error) {
	players := []string{}
	for p, _ := range f.playerAnswers {
		players = append(players, p)
	}
	return players, nil
}

// Get implements QuestionsGetter.
func (*fakeServerer) Get(context.Context) ([][]string, error) {
	return [][]string{{"Question?", "Answer."}}, nil
}

// Broadcast implements QuestionBroadcaster.
func (f *fakeServerer) Broadcast(context.Context, string, int, time.Duration, bool) (map[string]string, error) {
	return f.playerAnswers, nil
}

// Broadcast implements ResultsBroadcaster.
func (*fakeServerer2) Broadcast(context.Context, map[string]int) error {
	return nil
}
