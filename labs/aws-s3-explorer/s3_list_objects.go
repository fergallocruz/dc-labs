package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	accessKey = "fer"
	secretKey = "12345"
	endpoint  = "https://s3.us-east-1.amazonaws.com"
	region    = "us-east-1"
	bucket    = ""
)

// Lists the items in the specified S3 Bucket
//
// Usage:
//    go run s3_list_objects.go BUCKET_NAME
func main() {
	if len(os.Args) != 2 {
		exitErrorf("Bucket name required\nUsage: %s bucket_name",
			os.Args[0])
	}

	bucket := os.Args[1]

	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(AnonymousCredentials),
		// switching between Endpoint and EndpointResolver reproduce the same issue
		Endpoint: aws.String(endpoint),
		// EndpointResolver: endpoints.ResolverFunc(resolver),
	}))

	client := s3.New(sess)
	o, err := client.ListObjects(&s3.ListObjectsInput{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int64(1),
	})

	fmt.Println(o, err)

	for _, item := range o.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}

	fmt.Println("Found", len(o.Contents), "items in bucket", bucket)
	fmt.Println("")
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
