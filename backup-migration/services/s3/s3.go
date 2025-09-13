package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"backup-migration/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	s3     *s3.Client
	bucket string
	prefix string
}

type Object struct {
	Key          string
	Size         int64
	LastModified time.Time
	Generation   string
	SnapshotNum  int64
}

func NewClient(cfg *models.Configs) (*Client, error) {
	if cfg == nil || cfg.S3_secrets_isMissing {
		return nil, errors.New("S3 Initialized without config")
	}

	s3Config, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.S3_secrets.AWS_ACCESS_KEY_ID, cfg.S3_secrets.AWS_SECRET_ACCESS_KEY, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("couldn't load default configuration. %w", err)
	}

	client := s3.NewFromConfig(s3Config, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(cfg.S3_secrets.AWS_ENDPOINT_URL_S3)
		options.Region = cfg.S3_secrets.AWS_REGION
	})

	return &Client{
		s3:     client,
		prefix: cfg.S3_secrets.DB_PREFIX,
		bucket: cfg.S3_secrets.BUCKET_NAME,
	}, nil
}

// Browse lists snapshot objects under your structure.
func (c *Client) Browse(ctx context.Context, prefix string, max int32) ([]Object, error) {
	genID := "" // which generation to list

	// Normalize the user's prefix to a generation ID if provided.
	p := strings.TrimPrefix(prefix, "/")
	if p != "" {
		// Accept "generations/<gen>" or ".../generations/<gen>/snapshots"
		parts := strings.Split(p, "/")
		for i := 0; i < len(parts); i++ {
			if parts[i] == "generations" && i+1 < len(parts) {
				genID = parts[i+1]
				break
			}
		}
	}

	if genID == "" {
		// Find the latest generation by inspecting each generation's snapshots
		gens, err := c.listGenerations(ctx)
		if err != nil {
			return nil, err
		}
		if len(gens) == 0 {
			return nil, nil
		}
		// Choose the generation with the highest last snapshot number (tie-break: newest LastModified)
		var bestGen string
		var bestNum int64 = -1
		var bestTime time.Time

		for _, g := range gens {
			snaps, err := c.listSnapshots(ctx, g, 0) // 0 -> no per-page max; we sort/limit later
			if err != nil {
				return nil, fmt.Errorf("list snapshots for %s: %w", g, err)
			}
			if len(snaps) == 0 {
				continue
			}
			// snapshots already parsed with SnapshotNum; find max for this generation
			sort.Slice(snaps, func(i, j int) bool { return snaps[i].SnapshotNum > snaps[j].SnapshotNum })
			top := snaps[0]
			if top.SnapshotNum > bestNum || (top.SnapshotNum == bestNum && top.LastModified.After(bestTime)) {
				bestNum, bestTime, bestGen = top.SnapshotNum, top.LastModified, g
			}
		}
		if bestGen == "" {
			return nil, nil
		}
		genID = bestGen
	}

	// Return that generation's snapshots, sorted desc, limited by max.
	snaps, err := c.listSnapshots(ctx, genID, 0)
	if err != nil {
		return nil, err
	}
	sort.Slice(snaps, func(i, j int) bool { return snaps[i].SnapshotNum > snaps[j].SnapshotNum })

	if max > 0 && int32(len(snaps)) > max {
		snaps = snaps[:max]
	}
	return snaps, nil
}

// Download fetches the S3 object 'key' and writes it into 'outDir'.
// Returns the full local path written.
func (c *Client) Download(ctx context.Context, key, outDir string) (string, error) {
	key = c.joinPrefix(key)
	getOut, err := c.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &c.bucket,
		Key:    &key,
	})
	if err != nil {
		return "", fmt.Errorf("get object %q: %w", key, err)
	}
	defer getOut.Body.Close()

	base := filepath.Base(key)
	localPath := filepath.Join(outDir, base)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", fmt.Errorf("mkdir %q: %w", outDir, err)
	}

	f, err := os.Create(localPath)
	if err != nil {
		return "", fmt.Errorf("create %q: %w", localPath, err)
	}
	defer f.Close()

	if _, err := io.Copy(f, getOut.Body); err != nil {
		return "", fmt.Errorf("write %q: %w", localPath, err)
	}
	return localPath, nil
}

// Upload streams a local file to S3 at 'key' (under client's prefix if set).
func (c *Client) Upload(ctx context.Context, key, localPath string) error {
	key = c.joinPrefix(key)

	f, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("open %q: %w", localPath, err)
	}
	defer f.Close()

	_, err = c.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &c.bucket,
		Key:    &key,
		Body:   f,
		// ContentType: aws.String(detectContentType(localPath)), // optional
	})
	if err != nil {
		return fmt.Errorf("put object %q: %w", key, err)
	}
	return nil
}

func (c *Client) GetExistingPrefixes(ctx context.Context) {
	start := c.joinPrefix("")
	fmt.Print(start)
}

func (c *Client) joinPrefix(key string) string {
	k := strings.TrimPrefix(key, "/")
	if c.prefix == "" {
		return k
	}
	return strings.TrimSuffix(c.prefix, "/") + "/" + k
}

// listGenerations returns generation IDs under "<prefix>/generations/"
func (c *Client) listGenerations(ctx context.Context) ([]string, error) {
	base := c.joinPrefix("generations/")
	input := &s3.ListObjectsV2Input{
		Bucket:    &c.bucket,
		Prefix:    &base,
		Delimiter: aws.String("/"), // important: get CommonPrefixes = "subfolders"
	}
	pager := s3.NewListObjectsV2Paginator(c.s3, input)

	gens := make([]string, 0, 8)
	for pager.HasMorePages() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("list generations: %w", err)
		}
		for _, cp := range page.CommonPrefixes {
			// cp.Prefix is something like "<prefix>/generations/<gen>/"
			p := aws.ToString(cp.Prefix)
			// extract the gen between the last two slashes
			// use path pkg (S3 keys are '/')
			// trim "<prefix>/generations/" + trailing "/"
			trimmed := strings.TrimSuffix(strings.TrimPrefix(p, base), "/")
			if trimmed != "" {
				gens = append(gens, trimmed)
			}
		}
	}
	return gens, nil
}

// listSnapshots lists *.snapshot.lz4 under "<prefix>/generations/<gen>/snapshots/"
func (c *Client) listSnapshots(ctx context.Context, gen string, max int32) ([]Object, error) {
	pref := c.joinPrefix(path.Join("generations", gen, "snapshots") + "/")
	in := &s3.ListObjectsV2Input{
		Bucket: &c.bucket,
		Prefix: &pref,
	}
	if max > 0 {
		in.MaxKeys = &max
	}
	pager := s3.NewListObjectsV2Paginator(c.s3, in)
	out := make([]Object, 0, 64)
	for pager.HasMorePages() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("list snapshots %s: %w", gen, err)
		}
		for _, it := range page.Contents {
			key := aws.ToString(it.Key)
			if !strings.HasSuffix(key, ".snapshot.lz4") {
				continue
			}
			num, ok := parseSnapshotNum(key)
			if !ok {
				continue
			}
			out = append(out, Object{
				Key:          key,
				Size:         *it.Size,
				LastModified: aws.ToTime(it.LastModified),
				Generation:   gen,
				SnapshotNum:  num,
			})
		}
	}
	return out, nil
}

// parseSnapshotNum extracts the <N> from ".../<N>.snapshot.lz4"
func parseSnapshotNum(key string) (int64, bool) {
	name := path.Base(key) // S3 keys use forward slashes
	name = strings.TrimSuffix(name, ".snapshot.lz4")
	if name == "" {
		return 0, false
	}
	n, err := strconv.ParseInt(name, 10, 64)
	if err != nil {
		return 0, false
	}
	return n, true
}
