package asana

import "net/url"

type Extractor interface {
	GetAllUsers() ([]User, error)
	GetAllProjects() ([]Project, error)
}

type extractor struct {
	apiclient APIClient
}

func NewExtractor(apiclient APIClient) Extractor {
	return &extractor{apiclient}
}

func (e extractor) defaultQuery() url.Values {
	var query url.Values = make(url.Values)
	query.Set("limit", "100")

	return query
}

func (e extractor) GetAllWorkspaces() ([]Workspace, error) {
	query := e.defaultQuery()
	workspaces, nextPage, err := e.apiclient.ListWorkspaces(query)
	if err != nil {
		return nil, err
	}

	if nextPage != nil {
		for {
			additionalWorkspaces, nextPage, err := e.apiclient.ListWorkspaces(query)
			if err != nil {
				return nil, err
			}

			workspaces = append(workspaces, additionalWorkspaces...)

			if nextPage == nil {
				break
			} else {
				query.Set("offset", nextPage.Offset)
			}
		}
	}

	return workspaces, nil
}

func (e extractor) GetAllUsers() ([]User, error) {
	workspaces, err := e.GetAllWorkspaces()
	if err != nil {
		return nil, err
	}

	query := e.defaultQuery()
	var usersRes []User = make([]User, 0, len(workspaces)*100*5)
	for _, ws := range workspaces {
		for {
			query.Set("workspace", ws.GID)
			users, nextPage, err := e.apiclient.ListUsers(query)
			if err != nil {
				return nil, err
			}
			usersRes = append(usersRes, users...)

			if nextPage == nil {
				break
			} else {
				query.Set("offset", nextPage.Offset)
			}
		}
	}

	return usersRes, nil
}

func (e extractor) GetAllProjects() ([]Project, error) {
	workspaces, err := e.GetAllWorkspaces()
	if err != nil {
		return nil, err
	}

	query := e.defaultQuery()
	projectsRes := make([]Project, 0, len(workspaces)*100*5)
	for _, ws := range workspaces {
		for {
			query.Set("workspace", ws.GID)
			projects, nextPage, err := e.apiclient.ListProjects(query)
			if err != nil {
				return nil, err
			}
			projectsRes = append(projectsRes, projects...)

			if nextPage == nil {
				break
			} else {
				query.Set("offset", nextPage.Offset)
			}
		}
	}

	return projectsRes, nil
}
