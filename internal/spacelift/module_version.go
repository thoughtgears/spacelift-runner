package spacelift

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
)

func (c *Client) GetLatestModuleVersion() (map[string]string, error) {
	ctx := context.TODO()

	var response struct {
		Modules []struct {
			Name   string `json:"name"`
			Latest struct {
				Number string `json:"number"`
			} `json:"latest"`
		} `json:"modules"`
	}

	query := `
		query Modules {
		  modules {
			name
			latest {
			  number
			}
		  }
		}
	`

	req := graphql.NewRequest(query)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.jwtToken))

	if err := c.graphQLClient.Run(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to get modules: %v", err)
	}

	modules := make(map[string]string)

	for _, module := range response.Modules {
		modules[module.Name] = module.Latest.Number
	}

	return modules, nil
}
