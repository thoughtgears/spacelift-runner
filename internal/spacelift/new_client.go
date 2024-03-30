package spacelift

import (
	"context"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/machinebox/graphql"
)

// Client is the Spacelift client that reads environment variables and retrieves a JWT token
// for the Spacelift API to be able to make queries to the spacelift API.
type Client struct {
	BaseURL       string
	KeyID         string `envconfig:"SPACELIFT_KEY_ID" required:"true"`
	KeySecret     string `envconfig:"SPACELIFT_KEY_SECRET" required:"true"`
	graphQLClient *graphql.Client
	jwtToken      string
}

func NewClient() (*Client, error) {
	baseURL := "https://checkatrade.app.spacelift.io/graphql"

	client := Client{
		BaseURL:       baseURL,
		graphQLClient: graphql.NewClient(baseURL),
	}

	if err := envconfig.Process("", &client); err != nil {
		return nil, fmt.Errorf("failed to read environment variables: %v", err)
	}

	if err := client.getJWTToken(); err != nil {
		return nil, err
	}

	return &client, nil
}

// getJWTToken retrieves a JWT token from the Spacelift API using the provided token
// and stores it in the client
func (c *Client) getJWTToken() error {
	var resp struct {
		ApiKeyUser struct {
			Jwt string
		}
	}

	mutation := `
        mutation GetSpaceliftToken($keyId: ID!, $keySecret: String!) {
            apiKeyUser(id: $keyId, secret: $keySecret) {
                jwt
            }
        }
	`

	variables := struct {
		KeyID     string `json:"keyId"`
		KeySecret string `json:"keySecret"`
	}{
		KeyID:     c.KeyID,
		KeySecret: c.KeySecret,
	}

	req := graphql.NewRequest(mutation)
	req.Var("keyId", variables.KeyID)
	req.Var("keySecret", variables.KeySecret)
	req.Header.Set("Content-Type", "application/json")

	if err := c.graphQLClient.Run(context.TODO(), req, &resp); err != nil {
		return fmt.Errorf("failed to get JWT token for client: %v", err)
	}

	c.jwtToken = resp.ApiKeyUser.Jwt

	return nil
}
