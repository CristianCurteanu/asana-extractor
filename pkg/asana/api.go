package asana

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/CristianCurteanu/angler"
	slumber "github.com/CristianCurteanu/slumber"
)

var (
	errToManyRequests = errors.New("too many requests, retry")
)

type APIClient interface {
	ListProjects(query url.Values) ([]Project, *NextPage, error)
	ListWorkspaces(query url.Values) ([]Workspace, *NextPage, error)
	ListUsers(query url.Values) ([]User, *NextPage, error)
}

type apiClient struct {
	host        string
	accessToken string
}

func NewAPIClient(host, accessToken string) APIClient {
	return &apiClient{host, accessToken}
}

func (c *apiClient) ListUsers(query url.Values) ([]User, *NextPage, error) {
	var errResp ErrorsResponse
	handleStatusError := func(message string) angler.StatusHandlerFunc {
		return handleErrorStatusWithResponse(&errResp, message)
	}

	users, err := slumber.Retry(func() (MultipleResponse[User], error) {

		url := fmt.Sprintf("%s/users?%s", c.host, query.Encode())
		users, err := angler.Fetch[MultipleResponse[User]](
			angler.WithURL(url),
			angler.WithHeader("Authorization", fmt.Sprintf("Bearer %s", c.accessToken)),
			angler.WithStatusHandler(http.StatusTooManyRequests, handleStatusTooManyRequests(50*time.Millisecond)),
			angler.WithStatusHandler(http.StatusBadRequest, handleStatusError("missing of malformed parameter")),
			angler.WithStatusHandler(http.StatusUnauthorized, handleStatusError("unauthorized")),
			angler.WithStatusHandler(http.StatusNotFound, handleStatusError("not found")),
			angler.WithStatusHandler(http.StatusInternalServerError, handleStatusError("internal error, try again later")),
		)
		if err != nil && errors.Is(err, errToManyRequests) {
			return users, err
		}
		return users, nil
	},
		slumber.WithRetryPolicy(slumber.ExponentialBackoff),
		slumber.WithDelay(50*time.Millisecond),
		slumber.WithRetries(5),
	)

	if err != nil {
		return nil, nil, err
	}
	if len(errResp.Errors) != 0 {
		return nil, nil, fmt.Errorf("bad HTTP response status response: %+v", errResp)
	}

	return users.Data, users.NextPage, nil
}

func (c *apiClient) ListWorkspaceUsers(workspaceId string, query url.Values) ([]User, *NextPage, error) {
	var errResp *ErrorsResponse
	handleStatusError := func(message string) angler.StatusHandlerFunc {
		return handleErrorStatusWithResponse(&errResp, message)
	}
	users, err := slumber.Retry(func() (MultipleResponse[User], error) {

		users, err := angler.Fetch[MultipleResponse[User]](
			angler.WithURL(fmt.Sprintf("%s/workspaces/%s/users?%s", c.host, workspaceId, query.Encode())),
			angler.WithHeader("Authorization", fmt.Sprintf("Bearer %s", c.accessToken)),
			angler.WithStatusHandler(http.StatusTooManyRequests, handleStatusTooManyRequests(50*time.Millisecond)),
			angler.WithStatusHandler(http.StatusBadRequest, handleStatusError("missing of malformed parameter")),
			angler.WithStatusHandler(http.StatusUnauthorized, handleStatusError("unauthorized")),
			angler.WithStatusHandler(http.StatusNotFound, handleStatusError("not found")),
			angler.WithStatusHandler(http.StatusInternalServerError, handleStatusError("internal error, try again later")),
		)
		if err != nil && errors.Is(err, errToManyRequests) {
			return users, err
		}
		return users, nil
	},
		slumber.WithRetryPolicy(slumber.ExponentialBackoff),
		slumber.WithDelay(50*time.Millisecond),
		slumber.WithRetries(5),
	)
	if err != nil {
		return nil, nil, err
	}
	if errResp != nil {
		return nil, nil, fmt.Errorf("bad HTTP response status response: %+v", errResp)
	}

	return users.Data, users.NextPage, nil
}

func (c *apiClient) ListWorkspaces(query url.Values) ([]Workspace, *NextPage, error) {
	var errResp *ErrorsResponse
	handleStatusError := func(message string) angler.StatusHandlerFunc {
		return handleErrorStatusWithResponse(&errResp, message)
	}

	workspaces, err := slumber.Retry(func() (MultipleResponse[Workspace], error) {

		workspaces, err := angler.Fetch[MultipleResponse[Workspace]](
			angler.WithURL(fmt.Sprintf("%s/workspaces?%s", c.host, query.Encode())),
			angler.WithHeader("Authorization", fmt.Sprintf("Bearer %s", c.accessToken)),
			angler.WithStatusHandler(http.StatusTooManyRequests, handleStatusTooManyRequests(50*time.Millisecond)),
			angler.WithStatusHandler(http.StatusBadRequest, handleStatusError("missing of malformed parameter")),
			angler.WithStatusHandler(http.StatusUnauthorized, handleStatusError("unauthorized")),
			angler.WithStatusHandler(http.StatusNotFound, handleStatusError("not found")),
			angler.WithStatusHandler(http.StatusInternalServerError, handleStatusError("internal error, try again later")),
		)
		if err != nil && errors.Is(err, errToManyRequests) {
			return workspaces, err
		}
		return workspaces, nil
	},
		slumber.WithRetryPolicy(slumber.ExponentialBackoff),
		slumber.WithDelay(50*time.Millisecond),
		slumber.WithRetries(5),
	)
	if err != nil {
		return nil, nil, err
	}
	if errResp != nil {
		return nil, nil, fmt.Errorf("bad HTTP response status response: %+v", errResp)
	}

	return workspaces.Data, workspaces.NextPage, nil
}

func (c *apiClient) ListProjects(query url.Values) ([]Project, *NextPage, error) {
	var errResp *ErrorsResponse
	handleStatusError := func(message string) angler.StatusHandlerFunc {
		return handleErrorStatusWithResponse(&errResp, message)
	}
	projects, err := slumber.Retry(func() (MultipleResponse[Project], error) {

		projects, err := angler.Fetch[MultipleResponse[Project]](
			angler.WithURL(fmt.Sprintf("%s/projects?%s", c.host, query.Encode())),
			angler.WithHeader("Authorization", fmt.Sprintf("Bearer %s", c.accessToken)),
			angler.WithStatusHandler(http.StatusTooManyRequests, handleStatusTooManyRequests(50*time.Millisecond)),
			angler.WithStatusHandler(http.StatusBadRequest, handleStatusError("missing of malformed parameter")),
			angler.WithStatusHandler(http.StatusUnauthorized, handleStatusError("unauthorized")),
			angler.WithStatusHandler(http.StatusNotFound, handleStatusError("not found")),
			angler.WithStatusHandler(http.StatusInternalServerError, handleStatusError("internal error, try again later")),
		)

		if err != nil && errors.Is(err, errToManyRequests) {
			return projects, err
		}
		return projects, nil
	},
		slumber.WithRetryPolicy(slumber.ExponentialBackoff),
		slumber.WithDelay(50*time.Millisecond),
		slumber.WithRetries(5),
	)
	if err != nil {
		return nil, nil, err
	}
	if errResp != nil {
		return nil, nil, fmt.Errorf("bad HTTP response status response: %+v", errResp)
	}

	return projects.Data, projects.NextPage, nil
}

func handleErrorStatusWithResponse[T any](errorResponse *T, message string) angler.StatusHandlerFunc {
	return func(r *http.Response) (any, error) {
		respBody, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(respBody, errorResponse)
		if err != nil {
			return nil, err
		}
		return errorResponse, errors.New(message)
	}
}

func handleStatusTooManyRequests(sleepTime time.Duration) angler.StatusHandlerFunc {
	return func(resp *http.Response) (any, error) {
		sleepIfRetryAfter(&resp.Header, sleepTime)

		return nil, errToManyRequests
	}
}

func sleepIfRetryAfter(headers *http.Header, defaultDuration time.Duration) {
	if retryAfter := headers.Get("Retry-After"); retryAfter != "" {
		delaySeconds, err := strconv.Atoi(retryAfter)
		if err == nil {
			time.Sleep((time.Duration(delaySeconds) * time.Second) - defaultDuration)
		}
	}
}
