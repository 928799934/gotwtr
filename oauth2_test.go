package gotwtr_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/928799934/gotwtr"
)

func Test_generateAppOnlyBearerToken(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx            context.Context
		client         *http.Client
		consumerKey    string
		consumerSecret string
	}
	tests := []struct {
		name string
		args args
		want struct {
			BearerToken    string
			ConsumerKey    string
			ConsumerSecret string
		}
		wantErr bool
	}{
		{
			name: "200 ok",
			args: args{
				ctx: context.Background(),
				client: mockHTTPClient(func(request *http.Request) *http.Response {
					data := `{
						"token_type": "bearer",
						"access_token": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
					}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(data)),
					}
				}),
				consumerKey:    "consumerKey",
				consumerSecret: "consumerSecret",
			},
			want: struct {
				BearerToken    string
				ConsumerKey    string
				ConsumerSecret string
			}{
				BearerToken:    "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
				ConsumerKey:    "consumerKey",
				ConsumerSecret: "consumerSecret",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := gotwtr.New(
				"key",
				gotwtr.WithHTTPClient(tt.args.client),
				gotwtr.WithConsumerKey(tt.args.consumerKey),
				gotwtr.WithConsumerSecret(tt.args.consumerSecret),
			)
			b, err := c.GenerateAppOnlyBearerToken(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.GenerateAppOnlyBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !b {
				t.Errorf("client.GenerateAppOnlyBearerToken() = %v, want %v", b, true)
				return
			}
			nowc := c.ExportClient()
			cstate := struct {
				BearerToken    string
				ConsumerKey    string
				ConsumerSecret string
			}{
				BearerToken:    nowc["bearerToken"],
				ConsumerKey:    nowc["consumerKey"],
				ConsumerSecret: nowc["consumerSecret"],
			}
			if diff := cmp.Diff(tt.want, cstate); diff != "" {
				t.Errorf("client.GenerateAppOnlyBearerToken() diff = %v", diff)
				return
			}
		})
	}
}
