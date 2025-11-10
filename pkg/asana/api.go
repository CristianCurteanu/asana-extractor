package asana

import (
	"errors"
	"fmt"
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
	users, err := slumber.Retry(func() (MultipleResponse[User], error) {

		users, err := angler.Fetch[MultipleResponse[User]](
			angler.WithURL(fmt.Sprintf("%s/users?%s", c.host, query.Encode())),
			angler.WithHeader("Authorization", fmt.Sprintf("Bearer %s", c.accessToken)),
			angler.WithStatusHandler(http.StatusTooManyRequests, func(resp *http.Response) (any, error) {
				sleepIfRetryAfter(&resp.Header, 50*time.Millisecond)

				return nil, errToManyRequests
			}),
			angler.WithStatusHandler(http.StatusBadRequest, func(resp *http.Response) (any, error) {
				return nil, errors.New("missing of malformed parameter")
			}),
			angler.WithStatusHandler(http.StatusUnauthorized, func(resp *http.Response) (any, error) {
				return nil, errors.New("unauthorized")
			}),
			angler.WithStatusHandler(http.StatusNotFound, func(resp *http.Response) (any, error) {
				return nil, errors.New("not found")
			}),
			angler.WithStatusHandler(http.StatusInternalServerError, func(resp *http.Response) (any, error) {
				return nil, errors.New("internal error, retry")
			}),
		)
		if err != nil && errors.Is(err, errToManyRequests) {
			return users, err
		}
		return users, err
	},
		slumber.WithRetryPolicy(slumber.ExponentialBackoff),
		slumber.WithDelay(50*time.Millisecond),
		slumber.WithRetries(5),
	)

	if err != nil {
		return nil, nil, err
	}

	return users.Data, users.NextPage, nil
}

func (c *apiClient) ListWorkspaceUsers(workspaceId string, query url.Values) ([]User, *NextPage, error) {
	users, err := slumber.Retry(func() (MultipleResponse[User], error) {

		users, err := angler.Fetch[MultipleResponse[User]](
			angler.WithURL(fmt.Sprintf("%s/workspaces/%s/users?%s", c.host, workspaceId, query.Encode())),
			angler.WithHeader("Authorization", fmt.Sprintf("Bearer %s", c.accessToken)),
			angler.WithStatusHandler(http.StatusTooManyRequests, func(resp *http.Response) (any, error) {
				sleepIfRetryAfter(&resp.Header, 50*time.Millisecond)

				return nil, errToManyRequests
			}),
			angler.WithStatusHandler(http.StatusBadRequest, func(resp *http.Response) (any, error) {
				return nil, errors.New("missing of malformed parameter")
			}),
			angler.WithStatusHandler(http.StatusUnauthorized, func(resp *http.Response) (any, error) {
				return nil, errors.New("unauthorized")
			}),
			angler.WithStatusHandler(http.StatusNotFound, func(resp *http.Response) (any, error) {
				return nil, errors.New("not found")
			}),
			angler.WithStatusHandler(http.StatusInternalServerError, func(resp *http.Response) (any, error) {
				return nil, errors.New("internal error, retry")
			}),
		)
		if err != nil && errors.Is(err, errToManyRequests) {
			return users, err
		}
		return users, err
	},
		slumber.WithRetryPolicy(slumber.ExponentialBackoff),
		slumber.WithDelay(50*time.Millisecond),
		slumber.WithRetries(5),
	)
	if err != nil {
		return nil, nil, err
	}

	return users.Data, users.NextPage, nil
}

func (c *apiClient) ListWorkspaces(query url.Values) ([]Workspace, *NextPage, error) {
	workspaces, err := slumber.Retry(func() (MultipleResponse[Workspace], error) {

		workspaces, err := angler.Fetch[MultipleResponse[Workspace]](
			angler.WithURL(fmt.Sprintf("%s/workspaces?%s", c.host, query.Encode())),
			angler.WithHeader("Authorization", fmt.Sprintf("Bearer %s", c.accessToken)),
			angler.WithStatusHandler(http.StatusTooManyRequests, func(resp *http.Response) (any, error) {
				sleepIfRetryAfter(&resp.Header, 50*time.Millisecond)

				return nil, errToManyRequests
			}),
			angler.WithStatusHandler(http.StatusBadRequest, func(resp *http.Response) (any, error) {
				return nil, errors.New("missing of malformed parameter")
			}),
			angler.WithStatusHandler(http.StatusUnauthorized, func(resp *http.Response) (any, error) {
				return nil, errors.New("unauthorized")
			}),
			angler.WithStatusHandler(http.StatusNotFound, func(resp *http.Response) (any, error) {
				return nil, errors.New("not found")
			}),
			angler.WithStatusHandler(http.StatusInternalServerError, func(resp *http.Response) (any, error) {
				return nil, errors.New("internal error, retry")
			}),
		)
		if err != nil && errors.Is(err, errToManyRequests) {
			return workspaces, err
		}
		return workspaces, err
	},
		slumber.WithRetryPolicy(slumber.ExponentialBackoff),
		slumber.WithDelay(50*time.Millisecond),
		slumber.WithRetries(5),
	)
	if err != nil {
		return nil, nil, err
	}

	return workspaces.Data, workspaces.NextPage, nil
}

func (c *apiClient) ListProjects(query url.Values) ([]Project, *NextPage, error) {
	projects, err := slumber.Retry(func() (MultipleResponse[Project], error) {

		projects, err := angler.Fetch[MultipleResponse[Project]](
			angler.WithURL(fmt.Sprintf("%s/projects?%s", c.host, query.Encode())),
			angler.WithHeader("Authorization", fmt.Sprintf("Bearer %s", c.accessToken)),
			angler.WithStatusHandler(http.StatusTooManyRequests, func(resp *http.Response) (any, error) {
				sleepIfRetryAfter(&resp.Header, 50*time.Millisecond)

				return nil, errToManyRequests
			}),
			angler.WithStatusHandler(http.StatusBadRequest, func(resp *http.Response) (any, error) {
				return nil, errors.New("missing of malformed parameter")
			}),
			angler.WithStatusHandler(http.StatusUnauthorized, func(resp *http.Response) (any, error) {
				return nil, errors.New("unauthorized")
			}),
			angler.WithStatusHandler(http.StatusNotFound, func(resp *http.Response) (any, error) {
				return nil, errors.New("not found")
			}),
			angler.WithStatusHandler(http.StatusInternalServerError, func(resp *http.Response) (any, error) {
				return nil, errors.New("internal error, retry")
			}),
		)

		if err != nil && errors.Is(err, errToManyRequests) {
			return projects, err
		}
		return projects, err
	},
		slumber.WithRetryPolicy(slumber.ExponentialBackoff),
		slumber.WithDelay(50*time.Millisecond),
		slumber.WithRetries(5),
	)
	if err != nil {
		return nil, nil, err
	}

	return projects.Data, projects.NextPage, nil
}

func sleepIfRetryAfter(headers *http.Header, defaultDuration time.Duration) {
	if retryAfter := headers.Get("Retry-After"); retryAfter != "" {
		delaySeconds, err := strconv.Atoi(retryAfter)
		if err == nil {
			time.Sleep((time.Duration(delaySeconds) * time.Second) - defaultDuration)
		}
	}
}
