package utils

/*
import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c *ReplicaClient) OpenLTXFile(ctx context.Context, level int, minTXID, maxTXID ltx.TXID, offset, size int64) (io.ReadCloser, error) {
	if err := c.Init(ctx); err != nil {
		return nil, err
	}

	var rangeStr string
	if size > 0 {
		rangeStr = fmt.Sprintf("bytes=%d-%d", offset, offset+size-1)
	} else {
		rangeStr = fmt.Sprintf("bytes=%d-", offset)
	}

	// Build the key from the file info
	filename := ltx.FormatFilename(minTXID, maxTXID)
	key := c.Path + "/" + fmt.Sprintf("%04x/%s", level, filename)
	out, err := c.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(key),
		Range:  aws.String(rangeStr),
	})
	if err != nil {
		if isNotExists(err) {
			return nil, os.ErrNotExist
		}
		return nil, fmt.Errorf("s3: get object %s: %w", key, err)
	}
	return out.Body, nil
}

// WriteLTXFile writes an LTX file to the replica.
func (c *ReplicaClient) WriteLTXFile(ctx context.Context, level int, minTXID, maxTXID ltx.TXID, r io.Reader) (*ltx.FileInfo, error) {
	if err := c.Init(ctx); err != nil {
		return nil, err
	}

	rc := internal.NewReadCounter(r)

	filename := ltx.FormatFilename(minTXID, maxTXID)
	key := c.Path + "/" + fmt.Sprintf("%04x/%s", level, filename)
	out, err := c.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(key),
		Body:   rc,
	})
	if err != nil {
		return nil, fmt.Errorf("s3: upload to %s: %w", key, err)
	}

	// Build file info from the uploaded file
	info := &ltx.FileInfo{
		Level:     level,
		MinTXID:   minTXID,
		MaxTXID:   maxTXID,
		Size:      rc.N(),
		CreatedAt: time.Now(),
	}

	internal.OperationTotalCounterVec.WithLabelValues(ReplicaClientType, "PUT").Inc()
	internal.OperationBytesCounterVec.WithLabelValues(ReplicaClientType, "PUT").Add(float64(rc.N()))

	// ETag indicates successful upload
	if out.ETag == nil {
		return nil, fmt.Errorf("s3: upload failed: no ETag returned")
	}

	return info, nil
}
*/
