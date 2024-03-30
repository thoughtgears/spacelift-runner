package backup

import (
	"context"
	"fmt"
	"os"
)

func (c *Client) Execute(stateFile string) error {
	file, err := os.ReadFile(stateFile)
	if err != nil {
		return fmt.Errorf("failed to read state file: %v", err)
	}

	wc := c.svc.Bucket(c.bucket).Object(c.object).NewWriter(context.TODO())
	wc.ContentType = "application/json"

	if _, err := wc.Write(file); err != nil {
		return fmt.Errorf("failed to write state file to GCS: %v", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	return nil
}
