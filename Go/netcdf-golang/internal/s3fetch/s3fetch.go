package s3fetch

import (
	"bufio"
	"bytes"
	"context"
	"log"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/gcsblob"
)

// NOTE
// The buck web ui is at: https://console.cloud.google.com/storage/browser/national-water-model?authuser=0&pli=1
// https://console.cloud.google.com/storage/browser/national-water-model
// gs://national-water-model
// this is a public bucket named:  national-water-model

func GetS3FP(oid string) ([]byte, error) {
	ctx := context.Background()
	// blob.OpenBucket creates a *blob.Bucket from a URL.
	// This URL will open the bucket "my-bucket" using default credentials.
	bucket, err := blob.OpenBucket(ctx, "gs://national-water-model")
	if err != nil {
		return nil, err
	}
	defer bucket.Close()

	// blobReader, err := bucket.NewReader(ctx, "nwm.20180917/analysis_assim/nwm.t00z.analysis_assim.channel_rt.tm00.conus.nc", nil)
	blobReader, err := bucket.NewReader(ctx, oid, nil)

	var b bytes.Buffer
	bw := bufio.NewWriter(&b)
	n, err := blobReader.WriteTo(bw)
	if err != nil {
		return nil, err
	}
	bw.Flush()
	log.Printf("WriteTo count: %d", n)

	blobReader.Close()
	return b.Bytes(), nil
}
