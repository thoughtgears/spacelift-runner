package backup

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

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

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	objectName := fmt.Sprintf("%s/%s/%s", stack, time.Now().Format("2006-04-02"), filePath)
	bucket := client.Bucket(bucketName)
	writer := bucket.Object(objectName).NewWriter(ctx)
	if _, err = io.Copy(writer, file); err != nil {
		return fmt.Errorf("error uploading file: %v", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("error closing writer: %v", err)
	}

	fmt.Printf("Backup of %s done: gs://%s/%s\n", filePath, bucketName, objectName)

	return nil
}
