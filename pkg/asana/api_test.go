package asana

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AsanaAPIClientTestSuite struct {
	suite.Suite

	apiclient APIClient
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(AsanaAPIClientTestSuite))
}

func (ts *AsanaAPIClientTestSuite) SetupTest() {
	ts.apiclient = NewAPIClient(
		"https://app.asana.com/api/1.0",
		"2/1211287002275341/1211897635868807:161d5f970d3728f0a6d8eeb2080c86bc",
	)
}

func (ts *AsanaAPIClientTestSuite) Test_ListUsers() {
	var query url.Values = make(url.Values)
	query.Set("limit", "100")

	wss, _, err := ts.apiclient.ListWorkspaces(query)
	ts.Require().NoError(err)

	var usersRes []User = make([]User, 0, len(wss)*100*5)
	for _, ws := range wss {
		for {
			query.Set("workspace", ws.GID)
			users, nextPage, err := ts.apiclient.ListUsers(query)
			ts.Require().NoError(err)
			ts.Require().NotEmpty(users)
			usersRes = append(usersRes, users...)

			if nextPage == nil {
				break
			} else {
				query.Set("offset", nextPage.Offset)
			}
		}
	}
	ts.Require().NotEmpty(usersRes)
}
