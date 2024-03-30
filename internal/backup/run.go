package backup

import (
	"context"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/storage"
)

func Run(bucketName, stack, filePath string) error {
	ctx := context.Background()
	// Create a GCS client
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating GCS client: %v", err)
	}

	defer client.Close()

	file, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}

	var object string

	if strings.Contains(filePath, "/") {
		parts := strings.Split(filePath, "/")
		object = fmt.Sprintf("%s/%s", stack, parts[len(parts)-1])
	} else {
		object = fmt.Sprintf("%s/%s", stack, filePath)
	}

	bucket := client.Bucket(bucketName)
	writer := bucket.Object(object).NewWriter(ctx)
	writer.ContentType = "application/json"

	if _, err := writer.Write(file); err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("error closing writer: %v", err)
	}

	fmt.Printf("Backup of %s done: gs://%s/%s\n", filePath, bucketName, object)

	return nil
}
