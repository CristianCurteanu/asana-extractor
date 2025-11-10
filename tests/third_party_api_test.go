package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/CristianCurteanu/angler"
	slumber "github.com/CristianCurteanu/slumber"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
)

type ThirdPartyAPIsTest struct {
	suite.Suite
}

var (
	errToManyRequests = errors.New("too many requests, retry")
)

type ProfileData struct {
	Id        int    `json:"id"`
	Username  string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	Company   string `json:"company"`
	Repos     int    `json:"public_repos"`
	Gists     int    `json:"public_gists"`
	Followers int    `json:"followers"`
	Following int    `json:"following"`
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ThirdPartyAPIsTest))
}

func (s *ThirdPartyAPIsTest) Test_NewRequest() {
	defer gock.Off()

	gock.New("https://(.*).com").
		Get("/user").
		Reply(429).
		BodyString("foo foo")

	gock.New("https://(.*).com").
		Get("/user").
		Reply(429).
		BodyString("foo foo")

	gock.New("https://(.*).com").
		Get("/user").
		Reply(200).
		BodyString(
			string(s.encodeJSON(&ProfileData{
				Username: "nerdchap",
			}),
			),
		)

	profile, err := slumber.Retry(func() (ProfileData, error) {
		profile, err := angler.Fetch[ProfileData](
			angler.WithURL("https://api.github.com/user"),
			angler.WithStatusHandler(http.StatusTooManyRequests, func(resp *http.Response) (any, error) {
				sleepIfRetryAfter(&resp.Header)

				return nil, errToManyRequests
			}),
			angler.WithStatusHandler(http.StatusInternalServerError, func(*http.Response) (any, error) {
				return nil, errors.New("failed response")
			}),
		)
		if err != nil && errors.Is(err, errToManyRequests) {
			return ProfileData{}, err
		}
		return profile, err
	},
		slumber.WithRetryPolicy(slumber.ExponentialBackoff),
		slumber.WithDelay(50*time.Millisecond),
		slumber.WithRetries(3),
	)

	s.Require().NoError(err)
	s.Require().NotZero(profile, fmt.Sprintf(">>> profile: %+v", profile))
}

func (s *ThirdPartyAPIsTest) encodeJSON(data any) []byte {
	encoded, err := json.Marshal(data)
	s.Require().NoError(err)

	return encoded
}

func sleepIfRetryAfter(headers *http.Header) {
	if retryAfter := headers.Get("Retry-After"); retryAfter != "" {
		delaySeconds, err := strconv.Atoi(retryAfter)
		if err == nil {
			time.Sleep((time.Duration(delaySeconds) * time.Second) - 50*time.Millisecond)
		}
		retryTime, err := time.Parse(time.RFC1123, retryAfter)
		if err == nil {
			time.Sleep(time.Until(retryTime) - 50*time.Millisecond)
		}
	}
}
