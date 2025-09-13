package s3

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Regncon/conorganizer/backup-migration/config"
	"github.com/Regncon/conorganizer/backup-migration/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Object struct {
	Key          string
	Size         int64
	LastModified time.Time
	Generation   string
	SnapshotNum  int64
}

type S3Client struct {
	Client *awsS3.Client
}

func NewS3Client() *S3Client {
	return &S3Client{}
}

func (c *S3Client) Connect(cfg *config.S3Config) error {
	fmt.Print("Connecting to S3")
	if cfg == nil {
		return errors.New("S3 Connect called without config")
	}

	// Collect credentials
	accessId := strings.TrimSpace(cfg.Access_id)
	accessKey := strings.TrimSpace(cfg.Access_key)
	endpoint := strings.TrimSpace(cfg.Endpoint)
	bucket := strings.TrimSpace(cfg.Bucket)
	region := strings.TrimSpace(cfg.Region)
	//prefix := strings.TrimSpace(cfg.S3.Prefix)

	// Check for missing required values
	if accessId == "" || accessKey == "" || endpoint == "" || bucket == "" {
		return fmt.Errorf("missing required S3 secrets (id:%t key:%t region:%t bucket:%t)",
			accessId != "", accessKey != "", region != "", bucket != "")

	}

	// Construct credentials
	credentials := credentials.NewStaticCredentialsProvider(accessId, accessKey, "")
	s3Config, err := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithCredentialsProvider(credentials))
	if err != nil {
		return fmt.Errorf("couldn't load default configuration. %w", err)
	}

	// Assign new client
	client := awsS3.NewFromConfig(s3Config, func(options *awsS3.Options) {
		options.BaseEndpoint = aws.String(endpoint)
		options.Region = region
	})

	c.Client = client
	fmt.Println("New S3 client connection established")

	return nil
}

func (c *S3Client) ListExistingPrefixes(cfg *config.Config) (*[]string, error) {
	if c.Client == nil {
		return nil, errors.New("GetLatestBackup called without a valid S3 client")
	}

	fmt.Println("fetching prefixes")

	// Create new context for queries
	ctx := context.TODO()

	// Fetch a list of latest generations
	prefixList, err := c.Client.ListObjectsV2(ctx, &awsS3.ListObjectsV2Input{
		Bucket:    aws.String(cfg.S3.Bucket),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return nil, fmt.Errorf("error listing S3 prefixes %w", err)
	}
	if len(prefixList.CommonPrefixes) == 0 {
		return nil, fmt.Errorf("no prefixes found")
	}

	var prefixes []string
	for _, prefix := range prefixList.CommonPrefixes {
		prefixes = append(prefixes, strings.TrimSuffix(*prefix.Prefix, "/"))
	}

	ctx.Done()

	return &prefixes, nil
}

func (c *S3Client) GetLatestBackup(cfg *config.Config) (*S3Object, error) {
	if c.Client == nil {
		return nil, errors.New("GetLatestBackup called without a valid S3 client")
	}

	// Create new context for queries
	ctx := context.TODO()

	// Fetch a list of latest generations
	genList, err := c.Client.ListObjectsV2(ctx, &awsS3.ListObjectsV2Input{
		Bucket:    aws.String(cfg.S3.Bucket),
		Prefix:    aws.String(cfg.S3.Prefix + "/generations/"),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return nil, fmt.Errorf("error listing S3 prefixes %w", err)
	}
	if len(genList.CommonPrefixes) == 0 {
		return nil, fmt.Errorf("no prefixes found")
	}

	// Track latest snapshot
	var latestKey string
	var latestTime time.Time

	// Look trough generation list and construct new queries with full path
	for _, generation := range genList.CommonPrefixes {
		snapshotPrefix := *generation.Prefix + "snapshots/"

		snapshots, err := c.Client.ListObjectsV2(ctx, &awsS3.ListObjectsV2Input{
			Bucket: aws.String(cfg.S3.Bucket),
			Prefix: aws.String(snapshotPrefix),
		})
		if err != nil {
			continue
		}

		// Filter snapshots for relevant data and update tracker
		for _, snapshot := range snapshots.Contents {
			if strings.HasSuffix(*snapshot.Key, ".snapshot.lz4") && snapshot.LastModified.After(latestTime) {
				latestKey = *snapshot.Key
				latestTime = *snapshot.LastModified
			}
		}

	}

	ctx.Done()

	return &S3Object{
		Key:          latestKey,
		LastModified: latestTime,
	}, nil
}

func (c *S3Client) Download(cfg *config.Config, key string) (*string, error) {
	if key == "" {
		return nil, errors.New("S3 Download must be called with a key")
	}

	// New context
	ctx := context.TODO()

	getOut, err := c.Client.GetObject(ctx, &awsS3.GetObjectInput{
		Bucket: aws.String(cfg.S3.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch snapshot %s: %w", key, err)
	}
	defer getOut.Body.Close()

	// decompress
	dbContent, err := utils.DecompressSnapshot(getOut.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress snapshot: %w", err)
	}

	// get working dir for saving file
	ex, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to attain working dir %w", err)
	}

	// construct filename
	lastModified := aws.ToTime(getOut.LastModified)
	fileTime := lastModified.Format("20060102_0304")
	fileDir := filepath.Dir(ex)
	fileName := "regncon_" + fileTime + ".db"

	// save file
	newFile, err := utils.CreateFile(dbContent, fileDir, fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	return newFile, nil
}
