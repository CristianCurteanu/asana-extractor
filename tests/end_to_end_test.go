package tests

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/CristianCurteanu/asana-extractor/pkg/asana"
	"github.com/CristianCurteanu/asana-extractor/pkg/storage"
	"github.com/google/uuid"
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

// TODO: Finish this test
func (ts *EndToEndTestSuit) Test_EndToEndExtraction_Success() {
	usersData, err := readFile(filepath.Join(ts.wd, "fixtures", "users.json"))
	ts.Require().NoError(err)
	gock.New("https://(.*).com").
		Get("/api/1.0/users").
		ParamPresent("limit").
		ParamPresent("workspace").
		Reply(200).
		BodyString(string(usersData))

	gock.New("https://(.*).com").
		Get("/api/1.0/workspaces").
		ParamPresent("limit").
		Reply(200).
		BodyString(string(ts.encodeJSON(
			asana.MultipleResponse[asana.Workspace]{
				Data: []asana.Workspace{
					{
						GID: uuid.NewString(),
					},
				},
			},
		)))

		// Get All the data
	users, err := ts.extractor.GetAllUsers()
	ts.Require().NoError(err)
	ts.Require().NotEmpty(users)

	// Store all the data
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
