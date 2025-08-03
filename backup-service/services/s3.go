package services

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/Regncon/conorganizer/backup-service/config"
)

// NewS3Client creates and configures an S3 client from the provided config.
func NewS3Client(cfg config.Config, logger *slog.Logger) (*s3.Client, error) {
	logger.Info("Initializing Tigris S3 client", slog.Group("tigris",
		slog.String("endpoint", cfg.AWS_ENDPOINT_URL_S3),
		slog.String("region", cfg.AWS_REGION),
		slog.String("bucket", cfg.BUCKET_NAME),
	))

	s3Config, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AWS_ACCESS_KEY_ID, cfg.AWS_SECRET_ACCESS_KEY, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("couldn't load default configuration. %w", err)
	}

	client := s3.NewFromConfig(s3Config, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(cfg.AWS_ENDPOINT_URL_S3)
		options.Region = cfg.AWS_REGION
	})
	logger.Info("S3 client initialized successfully")

	// Return S3 service client
	return client, nil
}

type SnapshotMeta struct {
	Key          string
	GenerationID string
	LastModified time.Time
}

// DownloadLatestSnapshot finds and downloads the most recent snapshot.lz4 file from the latest generation.
func DownloadLatestSnapshot(ctx context.Context, s3Client *s3.Client, bucket string, dbPrefix string) (string, error) {
	// Step 1: Find all generations
	genPrefix := strings.TrimSuffix(dbPrefix, "/") + "/generations/"
	genList, err := s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Prefix:    aws.String(genPrefix),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return "", fmt.Errorf("listing generations failed: %w", err)
	}
	if len(genList.CommonPrefixes) == 0 {
		return "", fmt.Errorf("no generations found under %s", genPrefix)
	}

	var latestKey string
	var latestTime time.Time

	// Step 2: Loop through generations and check for latest snapshot
	for _, gen := range genList.CommonPrefixes {
		snapshotPrefix := *gen.Prefix + "snapshots/"

		snapList, err := s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
			Prefix: aws.String(snapshotPrefix),
		})
		if err != nil {
			continue // skip broken or empty generations
		}

		for _, obj := range snapList.Contents {
			if strings.HasSuffix(*obj.Key, ".snapshot.lz4") && obj.LastModified.After(latestTime) {
				latestKey = *obj.Key
				latestTime = *obj.LastModified
			}
		}
	}

	if latestKey == "" {
		return "", fmt.Errorf("no .snapshot.lz4 files found across generations")
	}

	// Step 3: Download the snapshot content to a temporary file
	// todo: fix checksum warn
	getOut, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(latestKey),
	})
	if err != nil {
		return "", fmt.Errorf("failed to fetch snapshot %s: %w", latestKey, err)
	}
	defer getOut.Body.Close()

	tmpDir := "/mnt/regncon/backup/tmp"
	os.MkdirAll(tmpDir, 0o755)

	tmpFile, err := os.CreateTemp(tmpDir, "*.snapshot.lz4")
	if err != nil {
		return "", fmt.Errorf("could not create temporary file: %w", err)
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, getOut.Body); err != nil {
		return "", fmt.Errorf("failed to write snapshot to temp file: %w", err)
	}

	return tmpFile.Name(), nil
}
