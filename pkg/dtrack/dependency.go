package dtrack

import "strconv"

type Dependency struct {
	Project   Project   `json:"project"`
	Component Component `json:"component"`
}

func (c Client) GetDependenciesForProject(projectUUID string) ([]Dependency, error) {
	page := 1
	hasMorePages := true
	dependencies := make([]Dependency, 0)
	for hasMorePages {
		res, err := c.restClient.R().
			SetPathParams(map[string]string{
				"uuid": projectUUID,
			}).
			SetResult(make([]Dependency, 0)).
			SetQueryParams(map[string]string{
				"pageSize":   "100",
				"pageNumber": strconv.Itoa(page),
			}).
			Get("/api/v1/dependency/project/{uuid}")
		if err != nil {
			return nil, err
		}

		if err = c.checkResponse(res, 200); err != nil {
			return nil, err
		}

		dependenciesOnPage, ok := res.Result().(*[]Dependency)
		if !ok {
			return nil, ErrInvalidResponseType
		}

		dependencies = append(dependencies, *dependenciesOnPage...)

		totalCount, _ := strconv.Atoi(res.Header().Get("X-Total-Count"))

		hasMorePages = len(dependencies) < totalCount
		page++
	}

	return dependencies, nil
}
