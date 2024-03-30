package backup

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
)

type Client struct {
	svc    *storage.Client
	bucket string
	object string
}

func NewClient(bucket, object string) (*Client, error) {
	ctx := context.TODO()

	client := &Client{
		bucket: bucket,
		object: object,
	}

	// Check if the GOOGLE_OAUTH_ACCESS_TOKEN environment variable is set.
	if accessToken, exists := os.LookupEnv("GOOGLE_OAUTH_ACCESS_TOKEN"); exists {
		// Use the OAuth2 access token for authentication.
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
		clientOption := option.WithTokenSource(ts)
		c, err := storage.NewClient(ctx, clientOption)
		if err != nil {
			return nil, fmt.Errorf("failed to create client: %v", err)
		}

		client.svc = c
	} else {
		// Fall back to ADC.
		c, err := storage.NewClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create client: %v", err)
		}

		client.svc = c
	}

	return client, nil
}
