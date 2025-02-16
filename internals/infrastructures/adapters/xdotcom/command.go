package xdotcom

import (
	"context"
	"fmt"
	"github.com/Pr3c10us/boilerplate/internals/domains/xdotcom"
	"github.com/Pr3c10us/boilerplate/packages/configs"
	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"
	"net/http"
	"os"
)

type Repository struct {
	environmentVariables *configs.EnvironmentVariables
}

func NewXDotComRepository(environmentVariables *configs.EnvironmentVariables) xdotcom.Repository {
	return &Repository{environmentVariables: environmentVariables}
}

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func (repo *Repository) Tweet(tweet xdotcom.Tweet) (string, error) {
	// Check expected secrets are set in the environment variables
	accessToken := repo.environmentVariables.XDotCom.AccessKey
	accessSecret := repo.environmentVariables.XDotCom.AccessSecret

	client, err := newOAuth1Client(accessToken, accessSecret)
	if err != nil {
		return "", err
	}

	tweetId, err := tweet_g(client, tweet.Text, tweet.PreviousTweetID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return "", err
	}

	return tweetId, nil
}

func newOAuth1Client(accessToken, accessSecret string) (*gotwi.Client, error) {
	in := &gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:           accessToken,
		OAuthTokenSecret:     accessSecret,
	}

	return gotwi.NewClient(in)
}

func tweet_g(c *gotwi.Client, text, id string) (string, error) {
	var p *types.CreateInput
	if id != "" {
		p = &types.CreateInput{
			Text: gotwi.String(text),
			Reply: &types.CreateInputReply{
				InReplyToTweetID: id,
			},
		}
	} else {
		p = &types.CreateInput{
			Text: gotwi.String(text),
		}
	}

	res, err := managetweet.Create(context.Background(), c, p)
	if err != nil {
		return "", err
	}

	return gotwi.StringValue(res.Data.ID), nil
}
