package utils

import (
	"fmt"
	"io"

	"github.com/pierrec/lz4/v4"
)

func DecompressSnapshot(reader io.Reader) ([]byte, error) {
	decompressed := lz4.NewReader(reader)

	data, err := io.ReadAll(decompressed)
	if err != nil {
		return nil, fmt.Errorf("error decompressing data: %w", err)
	}
	return data, nil
}

func CompressDatabase(writer io.Writer) ([]byte, error) {
	var test []byte
	return test, nil
}
