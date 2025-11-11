package asana

type MultipleResponse[T any] struct {
	Data     []T       `json:"data"`
	NextPage *NextPage `json:"next_page,omitempty"`
}

type NextPage struct {
	Offset string `json:"offset"`
}

type User struct {
	GID        string    `json:"gid"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Workspaces []Compact `json:"workspaces"`
	Photo      *Photo    `json:"photo,omitempty"`
}

type Workspace struct {
	GID string `json:"gid"`
}

type Project struct {
	GID string `json:"gid"`
}

type Compact struct {
	GID          string `json:"gid"`
	ResourceType string `json:"resource_type"`
	Name         string `json:"name"`
}

type Photo struct {
	Small  string `json:"photo.image_27x27"`
	Medium string `json:"photo.image_128x128"`
	Huge   string `json:"photo.image_1024x1024"`
}

type ErrorsResponse struct {
	Errors []ErrorResponse `json:"errors"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Help    string `json:"help"`
	Phrase  string `json:"phrase"`
}
