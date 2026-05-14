package clientjar

import (
	"archive/zip"
	"bytes"
)

type ClientJar struct {
	rd *zip.Reader
}

func Open(zipBytes []byte) (*ClientJar, error) {
	rd, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return nil, err
	}
	return &ClientJar{rd: rd}, nil
}
