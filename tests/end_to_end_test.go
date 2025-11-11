package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/CristianCurteanu/asana-extractor/pkg/asana"
	"github.com/CristianCurteanu/asana-extractor/pkg/storage"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
)

type EndToEndTestSuit struct {
	suite.Suite

	asanaUrl  string
	apiclient asana.APIClient
	extractor asana.Extractor
	fs        storage.File
	wd        string
	outputDir string
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EndToEndTestSuit))
}

func (ts *EndToEndTestSuit) SetupTest() {
	ts.asanaUrl = "https://app.asana.com/api/1.0"
	ts.apiclient = asana.NewAPIClient(
		ts.asanaUrl, "",
	)

	ts.extractor = asana.NewExtractor(ts.apiclient)
	wd, err := os.Getwd()
	ts.Require().NoError(err)
	ts.wd = wd
	ts.outputDir = filepath.Join(wd, "output")
	ts.fs = storage.NewFile(ts.outputDir)
}

func (ts *EndToEndTestSuit) Test_ExtractUsers_Success() {
	users, err := ts.extractor.GetAllUsers()
	ts.Require().NoError(err)
	ts.Require().NotEmpty(users)
}

func (ts *EndToEndTestSuit) Test_ExtractProjects_Success() {
	projects, err := ts.extractor.GetAllProjects()
	ts.Require().NoError(err)
	ts.Require().NotEmpty(projects)
}

func (ts *EndToEndTestSuit) Test_APIClient_ErrorIfUnauthorizedStatus() {
	gock.New("https://app.asana.com").
		Get("/api/1.0/users").
		ParamPresent("limit").
		Reply(http.StatusUnauthorized).
		BodyString(string(ts.encodeJSON(asana.ErrorsResponse{
			Errors: []asana.ErrorResponse{
				{
					Message: "unauthorized",
				},
			},
		})))
	_, _, err := ts.apiclient.ListUsers(url.Values{
		"limit": []string{"100"},
	})
	ts.Require().Error(err)
	ts.Require().ErrorContains(err, "bad HTTP response status response")
}

func (ts *EndToEndTestSuit) Test_APIClient_ErrorIfBadRequestStatus() {
	gock.New("https://app.asana.com").
		Get("/api/1.0/users").
		ParamPresent("limit").
		Reply(http.StatusBadRequest).
		BodyString(string(ts.encodeJSON(asana.ErrorsResponse{
			Errors: []asana.ErrorResponse{
				{
					Message: "bad request",
				},
			},
		})))
	_, _, err := ts.apiclient.ListUsers(url.Values{
		"limit": []string{"100"},
	})
	ts.Require().Error(err)
	ts.Require().ErrorContains(err, "bad HTTP response status response")
}

// TODO: Finish this test
func (ts *EndToEndTestSuit) Test_EndToEndExtraction_Success() {
	// usersData, err := readFile(filepath.Join(ts.wd, "fixtures", "users_response.json"))
	// ts.Require().NoError(err)
	// gock.New("https://app.asana.com").
	// 	Get("/api/1.0/users").
	// 	ParamPresent("limit").
	// 	Reply(200).
	// 	BodyString(string(usersData))

	// Get All the data
	// users, err := ts.extractor.GetAllUsers()
	// users, _, err := ts.apiclient.ListUsers(url.Values{
	// 	"limit": []string{"100"},
	// })
	// ts.Require().NoError(err)
	// ts.Require().NotEmpty(users)

	// Store all the data

	// Check if the files are stored
}

func (s *EndToEndTestSuit) encodeJSON(data any) []byte {
	encoded, err := json.Marshal(data)
	s.Require().NoError(err)

	return encoded
}

func readFile(path string) ([]byte, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err

	}

	return byteValue, nil
}
