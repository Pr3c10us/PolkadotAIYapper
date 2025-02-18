package command

import (
	"github.com/Pr3c10us/boilerplate/internals/domains/embedding"
	"github.com/Pr3c10us/boilerplate/internals/domains/llm"
	"github.com/Pr3c10us/boilerplate/internals/domains/xdotcom"
	"reflect"
	"testing"
)

func TestTweet_convertToArray(t *testing.T) {
	type fields struct {
		llm       llm.Repository
		embedding embedding.Repository
		xdotcom   xdotcom.Repository
	}
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "ss",
			fields: fields{
				llm:       nil,
				embedding: nil,
				xdotcom:   nil,
			},
			args:    args{input: "```json\n[\n    \"(1/7) � Ever heard of futarchy? Zeitgeist is shaking up #Polkadot with a g\novernance model that ties decisions to real-world events. Curious? It might just\n change how we think about decision-making! �� #Polkadot\",\n    \"(2/7) Imagine if governance decisions were based not just on votes, but on\ntangible outcomes. Zeitgeist's futarchy does just that, aligning incentives for\nmore effective results. What does this mean for #Polkadot?\",\n    \"(3/7) � Let's break it down: In futarchy, participants bet on the outcome\nof proposals. This betting reveals insights about potential success or failure.\nHow does this translate to better governance?\",\n    \"(4/7) Essentially, predictions come from those who will win or lose based o\nn real-world results. It's like having skin in the game, ensuring decisions serv\ne the community well. Curious about its impact on #Polkadot?\",\n    \"(5/7) � By integrating this model, #Polkadot could see more strategic and\ntransparent decision-making. It's a blend of democratic principles with market e\nfficiency. But there are challenges too. Let's explore!\",\n    \"(6/7) Critics argue risks in prediction markets, but proponents highlight i\nncreased accountability and innovation within #Polkadot. Zeitgeist is a pioneer;\n will others follow? What do you think?\",\n    \"(7/7) � Futarchy could redefine governance. Zeitgeist is leading the charg\ne on #Polkadot! Share your thoughts or ask questions below. Dive into the future\n of governance! #Innovation #Web3 #Blockchain\"\n]\n```"},
			want:    []string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Tweet{
				llm:       tt.fields.llm,
				embedding: tt.fields.embedding,
				xdotcom:   tt.fields.xdotcom,
			}
			got, err := service.convertToArray(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToArray() got = %v, want %v", got, tt.want)
			}
		})
	}
}
